#!/usr/bin/env python3
"""
GitHub Secrets Upload Script
Automatically reads secrets from github-secrets.txt and uploads them to GitHub repository
"""

import os
import sys
import re
import json
import subprocess
import argparse
from typing import Dict, Optional

class GitHubSecretsUploader:
    def __init__(self, user: str, repo: str):
        self.user = user
        self.repo = repo
        self.repo_full = f"{user}/{repo}"
        self.secrets = {}
        
    def check_gh_cli(self) -> bool:
        """Check if GitHub CLI is installed and authenticated"""
        try:
            result = subprocess.run(['gh', '--version'], 
                                  capture_output=True, text=True, check=True)
            print("✅ GitHub CLI found")
            
            # Check authentication
            auth_result = subprocess.run(['gh', 'auth', 'status'], 
                                       capture_output=True, text=True)
            if auth_result.returncode == 0:
                print("✅ GitHub CLI authenticated")
                return True
            else:
                print("❌ GitHub CLI not authenticated. Run: gh auth login")
                return False
                
        except (subprocess.CalledProcessError, FileNotFoundError):
            print("❌ GitHub CLI not found. Please install it first:")
            print("   winget install GitHub.cli")
            return False
    
    def check_repo_access(self) -> bool:
        """Verify repository access"""
        try:
            result = subprocess.run(['gh', 'repo', 'view', self.repo_full], 
                                  capture_output=True, text=True, check=True)
            print(f"✅ Repository access confirmed: {self.repo_full}")
            return True
        except subprocess.CalledProcessError:
            print(f"❌ Cannot access repository: {self.repo_full}")
            print("   Make sure you have access and the repository exists")
            return False
    
    def parse_secrets_file(self, filepath: str = "github-secrets.txt") -> bool:
        """Parse the github-secrets.txt file and extract all secrets"""
        if not os.path.exists(filepath):
            print(f"❌ Secrets file not found: {filepath}")
            return False
        
        print(f"📖 Reading secrets from {filepath}")
        
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Parse the file format
            current_name = None
            current_value = ""
            in_multiline = False
            brace_count = 0
            
            for line in content.split('\n'):
                line = line.strip()
                
                # Skip comments and empty lines
                if not line or line.startswith('#'):
                    continue
                
                # Check for Name: pattern
                if line.startswith('Name:'):
                    # Save previous secret if exists
                    if current_name and current_value.strip():
                        self.secrets[current_name] = current_value.strip()
                    
                    current_name = line.replace('Name:', '').strip()
                    current_value = ""
                    in_multiline = False
                    brace_count = 0
                    continue
                
                # Check for Value: pattern
                if line.startswith('Value:'):
                    value_part = line.replace('Value:', '').strip()
                    
                    # Check if it's a JSON object (starts with {)
                    if value_part.startswith('{'):
                        in_multiline = True
                        current_value = value_part
                        brace_count = value_part.count('{') - value_part.count('}')
                    else:
                        current_value = value_part
                        if current_name:
                            self.secrets[current_name] = current_value
                            current_name = None
                            current_value = ""
                    continue
                
                # Handle multiline JSON values
                if in_multiline and current_name:
                    current_value += "\n" + line
                    brace_count += line.count('{') - line.count('}')
                    
                    # Check if JSON object is complete
                    if brace_count <= 0:
                        in_multiline = False
                        self.secrets[current_name] = current_value.strip()
                        current_name = None
                        current_value = ""
            
            # Save last secret if exists
            if current_name and current_value.strip():
                self.secrets[current_name] = current_value.strip()
            
            print(f"✅ Parsed {len(self.secrets)} secrets from file")
            
            # Show what we found
            print("\n📋 Found secrets:")
            for name in self.secrets.keys():
                value_preview = self.secrets[name][:50] + "..." if len(self.secrets[name]) > 50 else self.secrets[name]
                print(f"  - {name}: {value_preview}")
            
            return len(self.secrets) > 0
            
        except Exception as e:
            print(f"❌ Error parsing secrets file: {e}")
            return False
    
    def validate_azure_credentials(self) -> bool:
        """Validate that AZURE_CREDENTIALS is valid JSON"""
        if 'AZURE_CREDENTIALS' not in self.secrets:
            print("❌ AZURE_CREDENTIALS not found in secrets")
            return False
        
        try:
            creds = json.loads(self.secrets['AZURE_CREDENTIALS'])
            required_fields = ['clientId', 'clientSecret', 'subscriptionId', 'tenantId']
            
            for field in required_fields:
                if field not in creds:
                    print(f"❌ Missing required field in AZURE_CREDENTIALS: {field}")
                    return False
                if not creds[field] or creds[field].startswith('your-'):
                    print(f"❌ Invalid placeholder value for {field}")
                    return False
            
            print("✅ AZURE_CREDENTIALS validation passed")
            return True
            
        except json.JSONDecodeError as e:
            print(f"❌ AZURE_CREDENTIALS is not valid JSON: {e}")
            return False
    
    def upload_secret(self, name: str, value: str) -> bool:
        """Upload a single secret to GitHub"""
        try:
            # Use gh secret set command
            cmd = ['gh', 'secret', 'set', name, '--repo', self.repo_full]
            
            result = subprocess.run(cmd, input=value, text=True, 
                                  capture_output=True, check=True)
            
            print(f"✅ {name} uploaded successfully")
            return True
            
        except subprocess.CalledProcessError as e:
            print(f"❌ Failed to upload {name}: {e}")
            if e.stderr:
                print(f"   Error: {e.stderr}")
            return False
    
    def upload_all_secrets(self) -> bool:
        """Upload all parsed secrets to GitHub"""
        if not self.secrets:
            print("❌ No secrets to upload")
            return False
        
        print(f"\n🚀 Uploading {len(self.secrets)} secrets to {self.repo_full}...")
        
        success_count = 0
        for name, value in self.secrets.items():
            if self.upload_secret(name, value):
                success_count += 1
        
        print(f"\n📊 Upload complete: {success_count}/{len(self.secrets)} secrets uploaded")
        
        if success_count == len(self.secrets):
            print("🎉 All secrets uploaded successfully!")
            self.show_repo_secrets()
            return True
        else:
            print("⚠️  Some secrets failed to upload")
            return False
    
    def show_repo_secrets(self):
        """Show current repository secrets"""
        try:
            result = subprocess.run(['gh', 'secret', 'list', '--repo', self.repo_full], 
                                  capture_output=True, text=True, check=True)
            print(f"\n📋 Current secrets in {self.repo_full}:")
            print(result.stdout)
        except subprocess.CalledProcessError:
            print("⚠️  Could not list repository secrets")
    
    def run(self) -> bool:
        """Main execution flow"""
        print("🔐 GitHub Secrets Upload Tool")
        print(f"📦 Repository: {self.repo_full}")
        print("=" * 50)
        
        # Check prerequisites
        if not self.check_gh_cli():
            return False
        
        if not self.check_repo_access():
            return False
        
        # Parse secrets file
        if not self.parse_secrets_file():
            return False
        
        # Validate Azure credentials
        if not self.validate_azure_credentials():
            return False
        
        # Upload secrets
        return self.upload_all_secrets()

def main():
    parser = argparse.ArgumentParser(description='Upload GitHub secrets from github-secrets.txt')
    parser.add_argument('--user', required=True, help='GitHub username')
    parser.add_argument('--repo', required=True, help='Repository name')
    parser.add_argument('--file', default='github-secrets.txt', help='Secrets file path')
    
    args = parser.parse_args()
    
    uploader = GitHubSecretsUploader(args.user, args.repo)
    success = uploader.run()
    
    if success:
        print("\n✅ Ready for deployment! You can now push your code.")
        sys.exit(0)
    else:
        print("\n❌ Upload failed. Please check the errors above.")
        sys.exit(1)

if __name__ == "__main__":
    main()
