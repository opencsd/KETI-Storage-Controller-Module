src="csd-identifier"

echo scp -rp $src root@10.0.4.84:/root/workspace/keti/CSD-Identifier copying...
sshpass -p ketidbms! scp -rp -o ConnectTimeout=60 $src root@10.0.4.84:/root/workspace/keti/CSD-Identifier


