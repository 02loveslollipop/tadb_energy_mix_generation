#!/bin/bash

# Azure Container Apps Deployment Script
# Run this script to set up the initial Azure infrastructure

# Variables - Update these with your values
RESOURCE_GROUP="rg-tadb-api"
LOCATION="East US"
CONTAINER_APP_ENV="env-tadb-api"
CONTAINER_APP_NAME="app-tadb-api"
ACR_NAME="acrtadbapi" # Must be globally unique, only lowercase letters and numbers
LOG_ANALYTICS_WORKSPACE="law-tadb-api"

echo "🚀 Setting up Azure Container Apps infrastructure..."

# Login to Azure (if not already logged in)
echo "📝 Logging in to Azure..."
az login

# Create resource group
echo "📦 Creating resource group..."
az group create \
  --name $RESOURCE_GROUP \
  --location "$LOCATION"

# Create Log Analytics workspace
echo "📊 Creating Log Analytics workspace..."
az monitor log-analytics workspace create \
  --resource-group $RESOURCE_GROUP \
  --workspace-name $LOG_ANALYTICS_WORKSPACE \
  --location "$LOCATION"

# Get Log Analytics workspace ID and key
LOG_ANALYTICS_WORKSPACE_ID=$(az monitor log-analytics workspace show \
  --resource-group $RESOURCE_GROUP \
  --workspace-name $LOG_ANALYTICS_WORKSPACE \
  --query customerId \
  --output tsv)

LOG_ANALYTICS_WORKSPACE_KEY=$(az monitor log-analytics workspace get-shared-keys \
  --resource-group $RESOURCE_GROUP \
  --workspace-name $LOG_ANALYTICS_WORKSPACE \
  --query primarySharedKey \
  --output tsv)

# Create Container Apps environment
echo "🌍 Creating Container Apps environment..."
az containerapp env create \
  --name $CONTAINER_APP_ENV \
  --resource-group $RESOURCE_GROUP \
  --location "$LOCATION" \
  --logs-workspace-id $LOG_ANALYTICS_WORKSPACE_ID \
  --logs-workspace-key $LOG_ANALYTICS_WORKSPACE_KEY

# Create Azure Container Registry
echo "🏪 Creating Azure Container Registry..."
az acr create \
  --resource-group $RESOURCE_GROUP \
  --name $ACR_NAME \
  --sku Basic \
  --location "$LOCATION"

# Enable admin user for ACR (for GitHub Actions)
az acr update \
  --name $ACR_NAME \
  --admin-enabled true

# Get ACR credentials
ACR_USERNAME=$(az acr credential show \
  --name $ACR_NAME \
  --query username \
  --output tsv)

ACR_PASSWORD=$(az acr credential show \
  --name $ACR_NAME \
  --query passwords[0].value \
  --output tsv)

# Create Container App
echo "📱 Creating Container App..."
az containerapp create \
  --name $CONTAINER_APP_NAME \
  --resource-group $RESOURCE_GROUP \
  --environment $CONTAINER_APP_ENV \
  --image mcr.microsoft.com/azuredocs/containerapps-helloworld:latest \
  --target-port 8080 \
  --ingress external \
  --cpu 0.25 \
  --memory 0.5Gi \
  --min-replicas 0 \
  --max-replicas 3

echo "✅ Infrastructure setup complete!"
echo ""
echo "📋 GitHub Secrets to configure:"
echo "REGISTRY_NAME: $ACR_NAME"
echo "REGISTRY_USERNAME: $ACR_USERNAME"
echo "REGISTRY_PASSWORD: $ACR_PASSWORD"
echo "RESOURCE_GROUP: $RESOURCE_GROUP"
echo "CONTAINER_APP_NAME: $CONTAINER_APP_NAME"
echo "CONTAINER_APP_ENVIRONMENT: $CONTAINER_APP_ENV"
echo ""
echo "🔑 Additional secrets you'll need to configure:"
echo "AZURE_CREDENTIALS: (Service Principal JSON)"
echo "DB_HOST: (Your PostgreSQL server host)"
echo "DB_PORT: 5432"
echo "DB_USER: (Your database user)"
echo "DB_PASSWORD: (Your database password)"
echo "DB_NAME: (Your database name)"
echo ""
echo "🌐 Your app URL will be available at:"
az containerapp show \
  --name $CONTAINER_APP_NAME \
  --resource-group $RESOURCE_GROUP \
  --query properties.configuration.ingress.fqdn \
  --output tsv
