from __future__ import annotations
import sys
import os
import argparse
import textwrap as _textwrap
from argparse import RawTextHelpFormatter
import pyodbc
import json
import re
from enum import Enum, IntEnum, unique
import signal
import readline
import time
import math
from dataclasses import dataclass
from datetime import datetime, timedelta
from tabulate import tabulate

from concurrent import futures
import grpc
import wal_manager_pb2
import wal_manager_pb2_grpc

@dataclass
class Table :
    column:str = None
    datatype:str = None
    ctype:str = None
    clen:int = None
    precision:int = None

@dataclass
class Tables :
    table_name:str = None
    index_val:str = None
    req_falg:bool = None
    table_infos:list[Table] = None

@dataclass
class Wal :
    table_name:str = None
    index:str = None
    seq:int = None
    row:object = None

host_addr = ''
account = ''
pwd = ''
database = ''
cursor = None
tables = []

signal.signal(signal.SIGTSTP, signal.SIG_IGN)

print('\nwal-manager Ver 1.0.0 for Linux on x86_64 (GPL)')
print('Copyright (c) 2022, OPENSYSNET and/or its affiliates.\n')

def arg_parser():
    parser = argparse.ArgumentParser(add_help=False, prog='wal-manager', formatter_class=make_wide(argparse.RawTextHelpFormatter))
    parser.add_argument('--help', action='help',  help='Show this help message and exit')
    parser.add_argument('-t', '--type', type=str, default='mysql', required=False, help='DBMS Type (default: MySQL)')
    parser.add_argument('-h', '--host', type=str, default='localhost', required=False, help='DB Address (default: localhost)') 
    parser.add_argument('-u', '--user', type=str, required=True, help='Database Account')
    parser.add_argument('-p', '--pwd', type=str, required=True, help='Database Password')
    parser.add_argument('-d', '--database', type=str, required=True, help='Database Schema')
    return parser.parse_args()


def make_wide(formatter, w=200, h=100):
    try:
        kwargs = {'width': w, 'max_help_position': h}
        formatter(None, **kwargs)
        return lambda prog: formatter(prog, **kwargs)
    except TypeError:
        warnings.warn("argparse help formatter failed, falling back.")
        return formatter


def getTableInfo(table_name, table_infos) :
    global cursor
    table_query = 'SELECT  COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, CHARACTER_OCTET_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE FROM INFORMATION_SCHEMA.COLUMNS  WHERE TABLE_SCHEMA=\'' + database + '\'  AND TABLE_NAME= \'' + table_name + '\' order by ORDINAL_POSITION;'
    cursor.execute(table_query)
    query_res = cursor.fetchall()
    for query_item in query_res:
        table_item = Table()
        table_item.column = query_item[0]
        table_item.datatype = query_item[1]
        table_item.ctype = query_item[2]
        if table_item.datatype == 'char' or table_item.datatype == 'varchar' :
            table_item.clen = int(query_item[3])
        elif table_item.datatype == 'int' :
            table_item.clen = 4
        elif table_item.datatype == 'date' :
            table_item.clen = 3
        elif table_item.datatype == 'decimal' :
            table_item.clen = (int(query_item[4])-int(query_item[5]))//2 + int(query_item[5])//2
            table_item.precision = (int(query_item[4])-int(query_item[5]))//2
        table_infos.append(table_item)

def parsingWAL (tables) :
    print('[WAL PARSING]')
    log_file = ''
    log_ext = r'.log'
    log_dir = r'/usr/local/mysql/data/.rocksdb/'
    log_list = [os.path.join(log_dir,file) for file in os.listdir(log_dir) if file.endswith(log_ext)]
    print(log_list)
   
    wal_list = []
    with open(log_list[0], 'rb') as f:
        while True: 
            byte = f.read(1)
            if byte == b'' : 
                break
            else :
                f.seek(-1,1)

            wal_item = Wal()

            f.seek(4,1) #CRC
            hex_str = f.read(2).hex()
            log_len = ''.join([hex_str[i-2:i] for i in range(len(hex_str), 0, -2)])
            # print('log len (%d)' % int(log_len, 16))
            remain_len = int(log_len, 16)

            log_type = f.read(1).hex()
            # print('log type (%d)' % int(log_type, 16))

            hex_str = f.read(8).hex()
            log_seq = ''.join([hex_str[i-2:i] for i in range(len(hex_str), 0, -2)])
            wal_item.seq = int(log_seq, 16)
            # print('log seq (%d)' % int(log_seq, 16))
            remain_len = remain_len - 8

            hex_str = f.read(4).hex()
            log_count = ''.join([hex_str[i-2:i] for i in range(len(hex_str), 0, -2)])
            # print('log count (%d)' % int(log_count, 16))
            remain_len = remain_len - 4

            op_type = f.read(1).hex()
            # print('op type (%d)' % int(op_type, 16))
            remain_len = remain_len - 1

            act_type = f.read(1).hex()
            # print('act type (%d)' % int(act_type, 16))
            remain_len = remain_len - 1

            index_len = f.read(1).hex()
            # print('index len (%d)' % int(index_len, 16))
            remain_len = remain_len - 1

            index = f.read(int(index_len, 16)).hex()
            index_str = index[0:8] 
            wal_item.index = index
            # print('index [%s]' % index_str)
            remain_len = remain_len - int(index_len, 16)

            if int(op_type,16) == 0x09 :
                # print("BEGIN_PREPARE")
                hex_str = f.read(1).hex()
                data_len = int(hex_str,16)
                # print('data len (%d)' % data_len)
                remain_len = remain_len - 1
                hex_str = f.read(1).hex()
                tail_len = int(hex_str,16)
                remain_len = remain_len - 1
                data = f.read(data_len).hex()
                remain_len = remain_len - data_len
                #print(data)

                tbl = None
                for table_info in tables :
                    if table_info.index_val == int(index_str, 16) and table_info.req_flag == True:
                        tbl = table_info.table_infos
                        wal_item.table_name = table_info.table_name
                if tbl is None : 
                    f.seek(remain_len,1) 
                    continue

                if int(act_type, 16) == 0 : #DELTE
                    f.seek(remain_len,1) 
                    wal_list.append(wal_item)
                    continue

                #PUT process
                pos = 0
                new_row = []
                for tbl_item in tbl : 
                    read_len = tbl_item.clen * 2
                    if tbl_item.datatype == 'varchar' :
                        read_len = int(data[pos:pos+2], 16)*2
                        pos = pos + 2
                    hex_str = data[pos:pos+read_len]
                    if tbl_item.datatype == 'int' :
                        hex_str = ''.join([hex_str[i-2:i] for i in range(len(hex_str), 0, -2)])
                        new_row.append(int(hex_str,16))
                    elif tbl_item.datatype == 'decimal' : 
                        if int(hex_str[0],16) >= 8 :
                            tmp_list = list(hex_str)
                            tmp_list[0] = str(int(hex_str[0], 16)-8)
                            hex_str = ''.join(tmp_list)
                        decimal = hex_str[:(tbl_item.precision*2)]
                        frag = hex_str[(tbl_item.precision*2):]
                        fin_val = "{:.2f}".format(int(decimal,16) + int(frag,16)/100)
                        new_row.append(float(fin_val))
                    elif tbl_item.datatype == 'char' : 
                        new_row.append((bytes.fromhex(hex_str).decode('ASCII')).strip())
                    elif tbl_item.datatype == 'date' : 
                        hex_str = ''.join([hex_str[i-2:i] for i in range(len(hex_str), 0, -2)])
                        year = int(hex_str,16) // (16 * 32)
                        month = (int(hex_str,16) % ( 16 * 32 )) // 32
                        day = int(hex_str,16) - ((year * 16 * 32) +(month * 32))
                        new_row.append(str(year) + '-' + str(month).zfill(2) + '-' + str(day).zfill(2))
                    elif tbl_item.datatype == 'varchar' : 
                        new_row.append((bytes.fromhex(hex_str).decode('ASCII')).strip())

                    pos = pos + read_len

                wal_item.row = new_row
                # print(wal_item)
                wal_list.append(wal_item)

                f.seek(remain_len,1) 

            else :
                # print("OTHER")
                # print(remain_len)
                f.seek(remain_len,1)

    print('\n[WAL ROW (TOTAL: %d)]' % len(wal_list))
    print(wal_list)
    #Update and Delete Process 
    fin_wal_list = []
    for wal_item in wal_list :
        find_flag = False
        for fin_wal_item in fin_wal_list :
            if wal_item.index == fin_wal_item.index :
                find_flag = True
                if wal_item.seq > fin_wal_item.seq :
                    fin_wal_list.remove(fin_wal_item)
                    if wal_item.row is not None :
                        fin_wal_list.append(wal_item)
                    break
        if find_flag == False :
            fin_wal_list.append(wal_item)

    print('\n[FINAL WAL ROW (TOTAL: %d)]' % len(fin_wal_list))
    print(fin_wal_list)
    return fin_wal_list


def checkWalStatus(table_name):
    global cursor
    print("CHECKING WAL STATUS")
    wal_query = 'select SUM(f.NUM_ROWS) from information_schema.ROCKSDB_DDL d, information_schema.rocksdb_index_file_map f where d.index_number=f.index_number and d.table_schema = \'' + database + '\' and d.table_name= \''+ table_name + '\' and d.index_name = \'HIDDEN_PK_ID\''
    cursor.execute(wal_query)
    wal_res = cursor.fetchone()
    sst_rows = wal_res[0]
    print('  [SST NUM. ROWS] {:,}'.format(int(sst_rows)))
    wal_query = 'select COUNT(*) from ' + table_name
    cursor.execute(wal_query)
    wal_res = cursor.fetchone()
    db_rows = wal_res[0]
    print('  [ DB NUM. ROWS] {:,}'.format(int(db_rows)))
    if sst_rows != db_rows :
        print('   [WAL NOT SYNC]')
    else :
        print('   [WAL SYNC]')


class WalManager (wal_manager_pb2_grpc.WalManagerServicer):
    def process_query(self, request, context):
        global cursor
        global tables
        start = time.time()
        math.factorial(100000)
        table_list = []
        for tables_info in tables :
            tables_info.req_flag = False
        if request.type == 'QUERY' :
            print('[QUERY PARSING]')
            print('Query : ' + request.value)
            query = request.value.lower()
            p = re.compile('(?<=from)(.*?)(?=where)')
            p_table_list = p.findall(query)
            table_names = [name.strip() for name in p_table_list]
            table_list = []
            for table_name in table_names :
                tbl_split = table_name.split(',')
                if len(tbl_split) > 0 :
                    for tbl_split_item in tbl_split :
                        table_list.append(tbl_split_item.strip())
                else :
                    table_list.append(table_name.strip())
            table_set = set(table_list)
            table_list = list(table_set)
        elif request.type == 'TABLE' :
            tbl_split = request.value.split(',')
            for tbl_split_item in tbl_split :
                table_list.append(tbl_split_item.strip())
        print(table_list)
        for table_name in table_list :
            find_flag = False
            for tables_info in tables : #later check table modification
                if tables_info.table_name == table_name :
                    find_flag = True
                    tables_info.req_flag = True
                    break
            if find_flag == False :
               tables_info = Tables()
               tables_info.table_name = table_name
               table_infos = []
               getTableInfo(table_name, table_infos)
               tables_info.table_infos = table_infos
               tables_info.req_flag = True
               index_query = 'select distinct(d.index_number) from information_schema.ROCKSDB_DDL d, information_schema.rocksdb_index_file_map f where d.index_number=f.index_number and d.table_name = \'' + table_name + '\';'
               cursor.execute(index_query)
               index_res = cursor.fetchone()
               index_number = index_res[0]
               tables_info.index_val = index_number
               print(tables_info)
               #checkWalStatus(table_name)
               tables.append(tables_info)
        wal_res = parsingWAL(tables)

        proto_res = wal_manager_pb2.query_res()
        proto_res.key = request.key
        for wal_info in wal_res :
            wal_list_item = proto_res.wal_list.add()
            wal_list_item.table_name = wal_info.table_name
            wal_list_item.index_val = wal_info.index
            for tables_info in tables :
                if tables_info.table_name == wal_info.table_name :
                    index = 0
                    for table_item in tables_info.table_infos :
                        column_list_item = wal_list_item.column_list.add()
                        column_list_item.column = table_item.column
                        column_list_item.datatype = table_item.datatype
                        column_list_item.ctype = table_item.ctype
                        column_list_item.clen = table_item.clen
                        if table_item.precision is not None :
                            column_list_item.precision = table_item.precision
                        column_list_item.value = str(wal_info.row[index])
                        index = index + 1
        if len(wal_res) == 0 :
            wal_list_item = proto_res.wal_list.add()

        #print(proto_res)
        end = time.time()
        print(f"{end - start:.5f} sec")
        return proto_res


if __name__ == "__main__":
    args = arg_parser()

    host_addr = args.host
    account = args.user
    pwd = args.pwd 
    database = args.database
    grpc_server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    wal_manager_pb2_grpc.add_WalManagerServicer_to_server(WalManager(), grpc_server)
    grpc_server.add_insecure_port('localhost:50051')
    grpc_server.start()

    try : 
        db_conn_str = 'DRIVER={MySQL ODBC 5.3 Driver};SERVER=' + host_addr + ';DATABASE=' + database + ';UID=' + account + ';PWD=' + pwd
        cnxn = pyodbc.connect(db_conn_str)
        cursor = cnxn.cursor()

        query = ''
        query_str = ''

        while(True) :
            time.sleep(1)
    except pyodbc.Error as ex:
        print('[ERROR] : Access denied for user \'' + account + '\'@\'' + host_addr + '\' (using password: YES)\n')
