#!/usr/bin/env python3
"""
Azure Service Principal Setup Automation
Automatically creates Azure service principal and updates GitHub secrets

Requirements:
    pip install azure-identity azure-mgmt-authorization azure-mgmt-resource azure-cli-core

Usage:
    python setup-azure-credentials.py --subscription-id <sub-id> --user <github-user> --repo <github-repo>
"""

import argparse
import json
import subprocess
import sys
from datetime import datetime
from typing import Dict, Any

def run_command(command: str, check: bool = True) -> subprocess.CompletedProcess:
    """Run shell command and return result"""
    print(f"🔧 Running: {command}")
    result = subprocess.run(command, shell=True, capture_output=True, text=True)
    
    if check and result.returncode != 0:
        print(f"❌ Command failed: {command}")
        print(f"Error: {result.stderr}")
        sys.exit(1)
    
    return result

def check_azure_cli():
    """Check if Azure CLI is installed and user is logged in"""
    print("🔍 Checking Azure CLI...")
    
    # Check if az command exists
    result = run_command("az --version", check=False)
    if result.returncode != 0:
        print("❌ Azure CLI not found. Please install it first:")
        print("   https://docs.microsoft.com/en-us/cli/azure/install-azure-cli")
        sys.exit(1)
    
    # Check if logged in
    result = run_command("az account show", check=False)
    if result.returncode != 0:
        print("❌ Not logged into Azure CLI. Please run: az login")
        sys.exit(1)
    
    print("✅ Azure CLI ready")

def check_github_cli():
    """Check if GitHub CLI is installed and user is logged in"""
    print("🔍 Checking GitHub CLI...")
    
    # Check if gh command exists
    result = run_command("gh --version", check=False)
    if result.returncode != 0:
        print("❌ GitHub CLI not found. Please install it first:")
        print("   https://cli.github.com/")
        sys.exit(1)
    
    # Check if logged in
    result = run_command("gh auth status", check=False)
    if result.returncode != 0:
        print("❌ Not logged into GitHub CLI. Please run: gh auth login")
        sys.exit(1)
    
    print("✅ GitHub CLI ready")

def get_azure_account_info() -> Dict[str, str]:
    """Get current Azure account information"""
    print("📋 Getting Azure account information...")
    
    result = run_command('az account show --query "{subscriptionId: id, tenantId: tenantId, name: name}" --output json')
    account_info = json.loads(result.stdout)
    
    print(f"✅ Subscription: {account_info['name']} ({account_info['subscriptionId']})")
    print(f"✅ Tenant ID: {account_info['tenantId']}")
    
    return account_info

def create_service_principal(subscription_id: str, sp_name: str = None) -> Dict[str, Any]:
    """Create Azure service principal for GitHub Actions"""
    if not sp_name:
        timestamp = datetime.now().strftime("%Y%m%d-%H%M%S")
        sp_name = f"sp-github-tadb-api-{timestamp}"
    
    print(f"🔧 Creating service principal: {sp_name}")
    
    # Create service principal with contributor role
    command = (
        f'az ad sp create-for-rbac '
        f'--name "{sp_name}" '
        f'--role contributor '
        f'--scopes /subscriptions/{subscription_id} '
        f'--json-auth'
    )
    
    result = run_command(command)
    
    try:
        sp_info = json.loads(result.stdout)
        print("✅ Service principal created successfully")
        return sp_info
    except json.JSONDecodeError:
        print("❌ Failed to parse service principal output")
        print(f"Output: {result.stdout}")
        sys.exit(1)

def get_container_registry_info(resource_group: str = "rg-tadb-api") -> Dict[str, str]:
    """Get Azure Container Registry information"""
    print("🔍 Getting Container Registry information...")
    
    # List ACR in resource group
    result = run_command(f'az acr list --resource-group {resource_group} --query "[0]" --output json', check=False)
    
    if result.returncode == 0 and result.stdout.strip() != "null":
        acr_info = json.loads(result.stdout)
        registry_name = acr_info['name']
        login_server = acr_info['loginServer']
        
        print(f"✅ Found existing registry: {registry_name}")
        
        # Enable admin access
        run_command(f'az acr update --name {registry_name} --admin-enabled true')
        
        # Get credentials
        result = run_command(f'az acr credential show --name {registry_name} --output json')
        credentials = json.loads(result.stdout)
        
        return {
            "name": registry_name,
            "login_server": login_server,
            "username": credentials['username'],
            "password": credentials['passwords'][0]['value']
        }
    else:
        print("❌ No Container Registry found. Please create one first.")
        return None

def update_secrets_file(azure_credentials: Dict[str, Any], registry_info: Dict[str, str]):
    """Update the github-secrets.txt file with new credentials"""
    print("📝 Updating github-secrets.txt...")
    
    secrets_content = f"""# GitHub Secrets Configuration Template
# Generated on {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}

# === AZURE AUTHENTICATION ===
Name: AZURE_CREDENTIALS
Value: {json.dumps(azure_credentials, indent=2)}

# === CONTAINER REGISTRY ===
Name: REGISTRY_LOGIN_SERVER
Value: {registry_info['login_server']}

Name: REGISTRY_USERNAME  
Value: {registry_info['username']}

Name: REGISTRY_PASSWORD
Value: {registry_info['password']}

# === AZURE RESOURCES ===
Name: RESOURCE_GROUP
Value: rg-tadb-api

Name: CONTAINER_APP_NAME
Value: app-tadb-api

Name: CONTAINER_APP_ENVIRONMENT
Value: env-tadb-api

# === DATABASE ===
Name: DB_URI
Value: postgresql://development:npg_ynS73MgwTOUR@ep-young-bush-aedny251-pooler.c-2.us-east-2.aws.neon.tech/neondb?sslmode=require&channel_binding=require

Name: DB_NAME
Value: consumo_energetico"""

    with open('github-secrets.txt', 'w') as f:
        f.write(secrets_content)
    
    print("✅ github-secrets.txt updated")

def upload_secrets_to_github(github_user: str, github_repo: str):
    """Upload secrets to GitHub repository"""
    print(f"🚀 Uploading secrets to {github_user}/{github_repo}...")
    
    result = run_command(f'python upload-secrets.py --user {github_user} --repo {github_repo}')
    print("✅ Secrets uploaded to GitHub")

def main():
    parser = argparse.ArgumentParser(description='Setup Azure credentials for GitHub Actions')
    parser.add_argument('--subscription-id', help='Azure subscription ID (optional, will auto-detect)')
    parser.add_argument('--user', required=True, help='GitHub username')
    parser.add_argument('--repo', required=True, help='GitHub repository name')
    parser.add_argument('--sp-name', help='Service principal name (optional)')
    parser.add_argument('--resource-group', default='rg-tadb-api', help='Azure resource group name')
    
    args = parser.parse_args()
    
    print("🚀 Azure Service Principal Setup Automation")
    print("=" * 50)
    
    # Check prerequisites
    check_azure_cli()
    check_github_cli()
    
    # Get Azure account info
    account_info = get_azure_account_info()
    subscription_id = args.subscription_id or account_info['subscriptionId']
    
    # Create service principal
    azure_credentials = create_service_principal(subscription_id, args.sp_name)
    
    # Get container registry info
    registry_info = get_container_registry_info(args.resource_group)
    if not registry_info:
        sys.exit(1)
    
    # Update secrets file
    update_secrets_file(azure_credentials, registry_info)
    
    # Upload to GitHub
    upload_secrets_to_github(args.user, args.repo)
    
    print("\n🎉 Setup completed successfully!")
    print("\nNext steps:")
    print("1. Push your code to trigger the GitHub Actions workflow")
    print("2. Monitor the deployment at: https://github.com/{}/{}/actions".format(args.user, args.repo))
    
    # Show summary
    print(f"\n📋 Summary:")
    print(f"   Service Principal: {azure_credentials['clientId']}")
    print(f"   Container Registry: {registry_info['name']}")
    print(f"   Login Server: {registry_info['login_server']}")

if __name__ == "__main__":
    main()
