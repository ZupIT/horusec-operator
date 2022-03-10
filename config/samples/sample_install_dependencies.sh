#!/bin/bash
# Copyright 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

echo "Deleting cluster..."
kind delete cluster --name horusec

echo "Creating..."
kind create cluster --name horusec

echo "Adding bitnami..."
helm repo add bitnami https://charts.bitnami.com/bitnami

helm repo update

echo "Adding installing dependencies..."
helm install rabbitmq bitnami/rabbitmq
helm install postgresql --set auth.database=horusec_db bitnami/postgresql

echo "Getting user and pwd of dependencies..."
export POSTGRES_USERNAME="postgres"
export POSTGRES_PASSWORD=$(kubectl get secret postgresql -o jsonpath="{.data.postgres-password}" | base64 --decode)
export RABBITMQ_USERNAME="user"
export RABBITMQ_PASSWORD=$(kubectl get secret rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 --decode)
export JWT_SECRET="4ff42f67-5929-fc52-65f1-3afc77ad86d5"

# waits for postgres to be ready
kubectl wait --for=condition=ready pod postgresql-0 --timeout 300s

kubectl run postgresql-client --rm -it --restart='Never' --image docker.io/bitnami/postgresql --env="PGPASSWORD=$POSTGRES_PASSWORD" --command -- psql --host postgresql -U $POSTGRES_USERNAME -d horusec_db -p 5432 --no-password -c "create database horusec_analytic_db;"

echo "Creating secrets of dependencies..."
kubectl create secret generic horusec-platform-database --from-literal="username=$POSTGRES_USERNAME" --from-literal="password=$POSTGRES_PASSWORD"
kubectl create secret generic horusec-analytic-database --from-literal="username=$POSTGRES_USERNAME" --from-literal="password=$POSTGRES_PASSWORD"

kubectl create secret generic horusec-broker --from-literal="username=$RABBITMQ_USERNAME" --from-literal="password=$RABBITMQ_PASSWORD"

kubectl create secret generic horusec-jwt --from-literal=jwt-token=$JWT_SECRET
