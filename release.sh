echo "====>start build<===="
if [[ -e ./server ]]; then
  rm ./server
fi
CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build
echo "====>end build<===="

echo "====>start deploy<===="
Host="10.11.159.156"
User="root";
DestDir="/data_b/eward/leetcodeServer"
ServerProcessNumber=$(ssh ${User}@${Host} "ps -ef | grep ${DestDir}/server | grep -v grep | awk '{print \$2;}'");
DefaultSeverPort="7834"
DefaultMysqlPort="3306"
DefaultMysqlHost="10.11.159.156"
DefaultMysqlDatabase="1024";
DefaultMysqlUser="root";

read -p "please input mysql's host (default: ${DefaultMysqlHost}): " MysqlHost
if test -z "$MysqlHost";then
  MysqlHost=$DefaultMysqlHost
fi;
read -p "please input mysql's database (default: ${DefaultMysqlDatabase}): " MysqlDatabase
if test -z "$MysqlDatabase";then
  MysqlDatabase=$DefaultMysqlDatabase
fi;
read -p "please input mysql's user (default: ${DefaultMysqlUser}): " MysqlUser
if test -z "$MysqlUser"; then
  MysqlUser=$DefaultMysqlUser
fi
read -s -p "please input mysql's password: " MysqlPassword
echo ""
read -p "please input mysql's port (default: ${DefaultMysqlPort}): " MysqlPort
if test -z "$MysqlPort"; then
  MysqlPort=$DefaultMysqlPort;
fi;
read -p "please input server port (default: ${DefaultSeverPort}): " ListenPort
if test -z "$ListenPort"; then
  ListenPort=$DefaultSeverPort;
fi;
read -p "display your input info [y/n, default:n]: " ans
if test -z "$ans"; then
  ans="n"
elif [ $ans == "y" ]; then
  echo "MysqlUser=${MysqlUser} MysqlPassword=${MysqlPassword} MysqlPort=${MysqlPort} MysqlHost=${MysqlHost} MysqlDatabase=${MysqlDatabase} ListenPort=${ListenPort}"
fi;

if ((ServerProcessNumber)); then
#  kill
  echo "killing Server";
  ssh ${User}@${Host} "kill -9 ${ServerProcessNumber}";
fi;

echo "upload binary file"
scp -q server ${User}@${Host}:${DestDir};
echo "starting Server"
ssh ${User}@${Host} "IsProd=true MysqlUser=${MysqlUser} MysqlPassword=${MysqlPassword} MysqlPort=${MysqlPort} MysqlHost=${MysqlHost} MysqlDatabase=${MysqlDatabase} ListenPort=${ListenPort} nohup ${DestDir}/server >/dev/null 2>&1 &";
rm server
echo "====>end deploy<===="
