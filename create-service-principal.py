#!/usr/bin/env python3
"""
Simple Azure Service Principal Creator
Creates a new service principal and updates GitHub secrets automatically

Usage:
    python create-service-principal.py
"""

import json
import subprocess
import sys
from datetime import datetime

def run_az_command(command: str) -> dict:
    """Run Azure CLI command and return JSON result"""
    print(f"🔧 Running: az {command}")
    
    full_command = f"az {command}"
    result = subprocess.run(full_command, shell=True, capture_output=True, text=True)
    
    if result.returncode != 0:
        print(f"❌ Command failed: {result.stderr}")
        return None
    
    try:
        return json.loads(result.stdout)
    except json.JSONDecodeError:
        print(f"✅ Command output: {result.stdout}")
        return {"output": result.stdout}

def create_new_service_principal():
    """Create a new service principal with current timestamp"""
    print("🚀 Creating new Azure Service Principal...")
    
    # Get current account info
    account = run_az_command('account show --query "{subscriptionId: id, tenantId: tenantId, name: name}"')
    if not account:
        print("❌ Failed to get Azure account info. Make sure you're logged in: az login")
        return None
    
    print(f"📋 Subscription: {account['name']}")
    print(f"📋 Subscription ID: {account['subscriptionId']}")
    print(f"📋 Tenant ID: {account['tenantId']}")
    
    # Create service principal
    timestamp = datetime.now().strftime("%Y%m%d-%H%M%S")
    sp_name = f"sp-github-tadb-api-{timestamp}"
    
    sp_command = f"""ad sp create-for-rbac --name "{sp_name}" --role contributor --scopes /subscriptions/{account['subscriptionId']} --json-auth"""
    
    sp_result = run_az_command(sp_command)
    if not sp_result:
        return None
    
    print("✅ Service Principal created successfully!")
    print(f"   Client ID: {sp_result['clientId']}")
    
    return sp_result

def update_github_secrets_file(azure_credentials: dict):
    """Update the github-secrets.txt file with new Azure credentials"""
    print("📝 Updating github-secrets.txt with new credentials...")
    
    # Read current file
    try:
        with open('github-secrets.txt', 'r') as f:
            content = f.read()
    except FileNotFoundError:
        print("❌ github-secrets.txt not found")
        return False
    
    # Find and replace AZURE_CREDENTIALS section
    lines = content.split('\n')
    new_lines = []
    skip_until_next_section = False
    
    for line in lines:
        if line.startswith('Name: AZURE_CREDENTIALS'):
            # Replace the entire AZURE_CREDENTIALS section
            new_lines.append('Name: AZURE_CREDENTIALS')
            new_lines.append(f'Value: {json.dumps(azure_credentials, indent=2)}')
            skip_until_next_section = True
        elif skip_until_next_section and line.startswith('# ==='):
            # We've reached the next section
            skip_until_next_section = False
            new_lines.append('')
            new_lines.append(line)
        elif not skip_until_next_section:
            new_lines.append(line)
    
    # Write updated content
    with open('github-secrets.txt', 'w') as f:
        f.write('\n'.join(new_lines))
    
    print("✅ github-secrets.txt updated successfully")
    return True

def upload_to_github():
    """Upload updated secrets to GitHub"""
    print("🚀 Uploading updated secrets to GitHub...")
    
    result = subprocess.run(
        'python upload-secrets.py --user 02loveslollipop --repo api_matriz_enegertica_tadb',
        shell=True,
        capture_output=True,
        text=True
    )
    
    if result.returncode == 0:
        print("✅ Secrets uploaded successfully!")
        return True
    else:
        print(f"❌ Failed to upload secrets: {result.stderr}")
        return False

def main():
    print("🔐 Azure Service Principal Auto-Creator")
    print("=" * 40)
    
    # Create new service principal
    azure_creds = create_new_service_principal()
    if not azure_creds:
        sys.exit(1)
    
    # Update secrets file
    if not update_github_secrets_file(azure_creds):
        sys.exit(1)
    
    # Upload to GitHub
    if not upload_to_github():
        sys.exit(1)
    
    print("\n🎉 All done! New service principal created and uploaded to GitHub.")
    print("\nNew Service Principal Details:")
    print(f"   Client ID: {azure_creds['clientId']}")
    print(f"   Tenant ID: {azure_creds['tenantId']}")
    print(f"   Subscription ID: {azure_creds['subscriptionId']}")
    print("\n✅ Ready to test deployment!")

if __name__ == "__main__":
    main()
