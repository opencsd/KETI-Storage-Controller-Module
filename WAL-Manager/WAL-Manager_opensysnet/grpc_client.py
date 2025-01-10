from __future__ import print_function

import logging
import sys
import grpc
import wal_manager_pb2
import wal_manager_pb2_grpc


def run():
	# grpc의 불안정한 채널인 위에서 오픈한 50051 포트로 연결합니다.
    with grpc.insecure_channel('localhost:50051') as channel:
    	# stub을 생성해줍니다.
        stub = wal_manager_pb2_grpc.WalManagerStub(channel)
        
        # 요청을 보내고 결과를 받는데, 서버에서 지정한 메서드에 요청시 사용할 proto 메시지 형식으로 요청을 전송합니다.
        response = stub.process_query(wal_manager_pb2.query_req(key='0|0|nation', type='TABLE', value='nation'))
    # print("client received")
    # print(response)
    print("get request:")
    print(wal_manager_pb2.query_req(key='0|0|nation', type='TABLE', value='nation'))


if __name__ == '__main__':
    logging.basicConfig()
    run()


