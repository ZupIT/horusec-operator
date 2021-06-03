echo "Deleting cluster..."
kind delete cluster --name horusec

echo "Creating..."
kind create cluster --name horusec

echo "Adding bitnami..."
helm repo add bitnami https://charts.bitnami.com/bitnami

helm repo update

echo "Adding installing dependencies..."
helm install rabbitmq bitnami/rabbitmq
helm install postgresql --set postgresqlDatabase=horusec_db bitnami/postgresql

echo "Getting user and pwd of dependencies..."
export POSTGRES_USERNAME="postgres"
export POSTGRES_PASSWORD=$(kubectl get secret postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode)
export RABBITMQ_USERNAME="user"
export RABBITMQ_PASSWORD=$(kubectl get secret rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 --decode)
export JWT_SECRET="4ff42f67-5929-fc52-65f1-3afc77ad86d5"

echo "Creating secrets of dependencies..."
kubectl create secret generic database-username --from-literal=database-username=$POSTGRES_USERNAME
kubectl create secret generic database-password --from-literal=database-password=$POSTGRES_PASSWORD
kubectl create secret generic database-uri --from-literal=database-uri=postgresql://$POSTGRES_USERNAME:$POSTGRES_PASSWORD@postgresql:5432/horusec_db?sslmode=disable

kubectl create secret generic broker-username --from-literal=broker-username=$RABBITMQ_USERNAME
kubectl create secret generic broker-password --from-literal=broker-password=$RABBITMQ_PASSWORD

kubectl create secret generic jwt-token --from-literal=jwt-token=$JWT_SECRET

echo "Installing horusec-operator..."
make install

echo "FINISHED !!"

echo "Now up horusec-operator and apply your changes to up horusec services....."
echo "Ex.:"
echo "go run cmd/app/main.go"
echo "cd config/samples"
echo "kubectl apply -f install_v2alpha1_horusec.yaml"