TOKEN=token-01
CLUSTER_STATE=new
NAME_1=machine-1
NAME_2=machine-2
NAME_3=machine-3
HOST_1=127.0.0.1
HOST_2=127.0.0.1
HOST_3=127.0.0.1
CLUSTER=${NAME_1}=http://${HOST_1}:2380,${NAME_2}=http://${HOST_2}:2381,${NAME_3}=http://${HOST_3}:2382

# For machine 1
etcd --data-dir=data.etcd/1 --name ${NAME_1} \
	--initial-advertise-peer-urls http://${HOST_1}:2380 --listen-peer-urls http://${HOST_1}:2380 \
	--advertise-client-urls http://${HOST_1}:2379 --listen-client-urls http://${HOST_1}:2379 \
	--initial-cluster ${CLUSTER} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}

# For machine 2
etcd --data-dir=data.etcd/2 --name ${NAME_2} \
	--initial-advertise-peer-urls http://${HOST_2}:2381 --listen-peer-urls http://${HOST_2}:2381 \
	--advertise-client-urls http://${HOST_2}:2378 --listen-client-urls http://${HOST_2}:2378 \
	--initial-cluster ${CLUSTER} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}

# For machine 3
etcd --data-dir=data.etcd/3 --name ${NAME_3} \
	--initial-advertise-peer-urls http://${HOST_3}:2382 --listen-peer-urls http://${HOST_3}:2382 \
	--advertise-client-urls http://${HOST_3}:2377 --listen-client-urls http://${HOST_3}:2377 \
	--initial-cluster ${CLUSTER} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}