user=guest
pass=guest

/usr/sbin/rabbitmqadmin -u $user -p $pass declare queue name=incoming durable=true auto_delete=false
/usr/sbin/rabbitmqadmin -u $user -p $pass declare exchange name=incoming type=topic durable=true auto_delete=false
/usr/sbin/rabbitmqadmin -u $user -p $pass declare binding source=incoming destination=incoming destination_type=queue routing_key=""

/usr/sbin/rabbitmqadmin -u $user -p $pass declare queue name=jobs durable=true auto_delete=false
/usr/sbin/rabbitmqadmin -u $user -p $pass declare exchange name=jobs type=topic durable=true auto_delete=false
/usr/sbin/rabbitmqadmin -u $user -p $pass declare binding source=jobs destination=jobs destination_type=queue routing_key=""
