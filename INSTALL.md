# Requirement
## postgresql
```shell
mkdir software
cd software
# https://www.postgresql.org/ftp/source/v15.2/
wget https://ftp.postgresql.org/pub/source/v15.2/postgresql-15.2.tar.gz
tar -xf postgresql-15.2.tar.gz
cd postgresql-15.2
# centos 
yum install readline-devel.x86_64
yum install zlib-devel;
make
make install

ln -s /usr/local/pgsql/bin/pg_ctl ~/bin/pg_ctl
ln -s /usr/local/pgsql/bin/psql ~/bin/psql


mkdir .mydb && cd .mydb && mkdir pgsql-5432;
pg_ctl init -D pgsql-5432
pg_ctl -D pgsql-5432 -l logfile.log start
 
psql -d postgres

```
## redis

```shell
cd software
wget https://download.redis.io/redis-stable.tar.gz
tar -xf redis-stable.tar.gz
cd redis-stable 
make
make install
mv redis-stable /usr/local/redis
cd ~
redis-server --port 6379 --daemonize yes --requirepass 138678Mm
```

## pulsar
```shell
cd software
# https://pulsar.apache.org/download/
wget https://dlcdn.apache.org/pulsar/pulsar-3.0.0/apache-pulsar-3.0.0-bin.tar.gz
tar -xf apache-pulsar-3.0.0-bin.tar.gz
yum install java 
add comment for conf/pulsar_env.sh: # PULSAR_GC=${PULSAR_GC:-"-XX:+UseZGC -XX:+PerfDisableSharedMem -XX:+AlwaysPreTouch"}
https://www.oracle.com/java/technologies/downloads/

sudo ln -s /usr/lib/jvm/jdk-20/bin/java /usr/local/bin/java
./bin/pulsar-daemon start standalone
```

##java
```shell

wget https://download.oracle.com/java/20/latest/jdk-20_linux-x64_bin.rpm
rpm -ivh jdk-20_linux-x64_bin.rpm

```

```shell

mkdir /myresource/
chmod 0755  /myresource/
```
