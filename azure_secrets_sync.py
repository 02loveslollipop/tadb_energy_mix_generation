#!/usr/bin/env python3
"""
Azure + GitHub secrets sync

Creates or retrieves Azure credentials (Service Principal + ACR credentials)
and uploads them as GitHub repository secrets in one go.

Requirements:
  - Azure CLI (az) logged in: az login
  - GitHub CLI (gh) logged in: gh auth login

Usage examples:
  python azure_secrets_sync.py --user 02loveslollipop --repo api_matriz_enegertica_tadb \
    --resource-group rg-tadb-api --container-app-name app-tadb-api \
    --container-app-env env-tadb-api

  # If you already have an ACR name:
  python azure_secrets_sync.py --user <u> --repo <r> --acr-name <acrname>

Optional:
  --db-uri <postgres_uri> --db-name <name>
"""

import argparse
import json
import subprocess
import sys
from datetime import datetime
from typing import Dict, Optional
from urllib.parse import urlparse, parse_qs


def run(cmd: str, check: bool = True) -> subprocess.CompletedProcess:
    print(f"$ {cmd}")
    cp = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    if check and cp.returncode != 0:
        print(cp.stdout)
        print(cp.stderr)
        sys.exit(cp.returncode)
    return cp


def ensure_cli(require_az: bool = True, require_gh: bool = True):
    if require_az:
        if run("az --version", check=False).returncode != 0:
            print("Azure CLI not found. Install: https://aka.ms/azure-cli")
            sys.exit(1)
        if run("az account show", check=False).returncode != 0:
            print("Not logged into Azure. Run: az login")
            sys.exit(1)
    if require_gh:
        if run("gh --version", check=False).returncode != 0:
            print("GitHub CLI not found. Install: https://cli.github.com/")
            sys.exit(1)
        if run("gh auth status", check=False).returncode != 0:
            print("Not logged into GitHub. Run: gh auth login")
            sys.exit(1)


def get_subscription() -> Dict[str, str]:
    res = run('az account show --query "{subscriptionId:id, tenantId:tenantId, name:name}" -o json')
    return json.loads(res.stdout)


def parse_db_uri(db_uri: str) -> Optional[Dict[str, str]]:
    """Parse a PostgreSQL DB URI into discrete components for fallback secrets."""
    try:
        p = urlparse(db_uri)
        if p.scheme not in ("postgres", "postgresql"):
            return None
        user = p.username or ""
        pwd = p.password or ""
        host = p.hostname or ""
        port = str(p.port or 5432)
        db = (p.path or "/").lstrip("/")
        q = parse_qs(p.query or "")
        sslmode = (q.get("sslmode", [""]) or [""])[0]
        return {
            "DB_HOST": host,
            "DB_PORT": port,
            "DB_USER": user,
            "DB_PASSWORD": pwd,
            "DB_NAME": db,
            "DB_SSL_MODE": sslmode or "",
        }
    except Exception:
        return None


def ensure_sp(subscription_id: str, name_prefix: str = "sp-github-tadb-api") -> Dict[str, str]:
    stamp = datetime.utcnow().strftime("%Y%m%d-%H%M%S")
    sp_name = f"{name_prefix}-{stamp}"
    cmd = (
        f'az ad sp create-for-rbac --name "{sp_name}" --role contributor '
        f'--scopes /subscriptions/{subscription_id} --json-auth'
    )
    res = run(cmd)
    return json.loads(res.stdout)


def get_acr_info(resource_group: Optional[str], acr_name: Optional[str]) -> Optional[Dict[str, str]]:
    if acr_name:
        ls = run(f'az acr show --name {acr_name} -o json', check=False)
        if ls.returncode != 0:
            print(f"ACR {acr_name} not found")
            return None
        info = json.loads(ls.stdout)
    else:
        if not resource_group:
            print("Provide --resource-group when --acr-name is omitted")
            return None
        ls = run(f'az acr list --resource-group {resource_group} -o json', check=False)
        if ls.returncode != 0:
            print("Failed to list ACRs")
            return None
        arr = json.loads(ls.stdout)
        if not arr:
            print("No ACR found in resource group. Create one first.")
            return None
        info = arr[0]

    # Enable admin and fetch credentials
    run(f'az acr update --name {info["name"]} --admin-enabled true')
    cred = run(f'az acr credential show --name {info["name"]} -o json')
    c = json.loads(cred.stdout)
    return {
        "name": info["name"],
        "login_server": info["loginServer"],
        "username": c["username"],
        "password": c["passwords"][0]["value"],
    }


def gh_set(repo: str, name: str, value: str):
    print(f"Setting secret {name}")
    cp = subprocess.run(['gh', 'secret', 'set', name, '--repo', repo], input=value, text=True)
    if cp.returncode != 0:
        print(f"Failed to set {name}")
        sys.exit(cp.returncode)


def main():
    ap = argparse.ArgumentParser(description='Create/Retrieve Azure creds and upload to GitHub secrets')
    ap.add_argument('--user', required=True, help='GitHub username or org')
    ap.add_argument('--repo', required=True, help='GitHub repository name')
    ap.add_argument('--resource-group', help='Azure resource group (for ACR discovery)')
    ap.add_argument('--acr-name', help='Azure Container Registry name (optional)')
    ap.add_argument('--container-app-name', help='Container App name (optional)')
    ap.add_argument('--container-app-env', help='Container Apps environment name (optional)')
    ap.add_argument('--db-uri', help='Database connection URI (optional)')
    ap.add_argument('--db-name', help='Database name (optional)')
    ap.add_argument('--creds-file', help='Path to existing AZURE_CREDENTIALS JSON file (optional)')
    ap.add_argument('--creds-json', help='Inline AZURE_CREDENTIALS JSON string (optional)')
    ap.add_argument('--dry-run', action='store_true', help='Print what would be set without calling az/gh')
    ap.add_argument('--upload-db-components', action='store_true', help='Also upload DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME derived from --db-uri')
    args = ap.parse_args()

    repo_full = f"{args.user}/{args.repo}"

    if args.dry_run:
        print("[DRY RUN] Would create Service Principal and set: AZURE_CREDENTIALS")
        print("[DRY RUN] Would set registry secrets if ACR is available: REGISTRY_LOGIN_SERVER, REGISTRY_USERNAME, REGISTRY_PASSWORD")
        if args.resource_group:
            print(f"[DRY RUN] Would set RESOURCE_GROUP={args.resource_group}")
        if args.container_app_name:
            print(f"[DRY RUN] Would set CONTAINER_APP_NAME={args.container_app_name}")
        if args.container_app_env:
            print(f"[DRY RUN] Would set CONTAINER_APP_ENVIRONMENT={args.container_app_env}")
        if args.db_uri:
            print(f"[DRY RUN] Would set DB_URI (redacted)")
            if args.upload_db_components:
                parts = parse_db_uri(args.db_uri)
                if parts:
                    red = parts.copy(); red['DB_PASSWORD'] = '<redacted>'
                    print(f"[DRY RUN] Would set discrete DB_*: {json.dumps(red)}")
        if args.db_name:
            print(f"[DRY RUN] Would set DB_NAME={args.db_name}")
        return

    # Real execution
    require_az = True
    if args.creds_file or args.creds_json:
        # We can skip AZ login if creds provided AND no ACR lookup requested
        require_az = bool(args.acr_name or args.resource_group)
    ensure_cli(require_az=require_az, require_gh=True)

    # Resolve credentials
    if args.creds_json:
        sp = json.loads(args.creds_json)
    elif args.creds_file:
        with open(args.creds_file, 'r', encoding='utf-8') as fh:
            sp = json.load(fh)
    else:
        sub = get_subscription()
        sp = ensure_sp(sub['subscriptionId'])

    acr = None
    if args.acr_name or args.resource_group:
        acr = get_acr_info(args.resource_group, args.acr_name)

    # Upload secrets
    gh_set(repo_full, 'AZURE_CREDENTIALS', json.dumps(sp))
    if acr:
        gh_set(repo_full, 'REGISTRY_LOGIN_SERVER', acr['login_server'])
        gh_set(repo_full, 'REGISTRY_USERNAME', acr['username'])
        gh_set(repo_full, 'REGISTRY_PASSWORD', acr['password'])

    if args.resource_group:
        gh_set(repo_full, 'RESOURCE_GROUP', args.resource_group)
    if args.container_app_name:
        gh_set(repo_full, 'CONTAINER_APP_NAME', args.container_app_name)
    if args.container_app_env:
        gh_set(repo_full, 'CONTAINER_APP_ENVIRONMENT', args.container_app_env)
    if args.db_uri:
        gh_set(repo_full, 'DB_URI', args.db_uri)
        if args.upload_db_components:
            parts = parse_db_uri(args.db_uri)
            if parts:
                gh_set(repo_full, 'DB_HOST', parts['DB_HOST'])
                gh_set(repo_full, 'DB_PORT', parts['DB_PORT'])
                gh_set(repo_full, 'DB_USER', parts['DB_USER'])
                gh_set(repo_full, 'DB_PASSWORD', parts['DB_PASSWORD'])
                gh_set(repo_full, 'DB_NAME', parts['DB_NAME'])
                if parts['DB_SSL_MODE']:
                    gh_set(repo_full, 'DB_SSL_MODE', parts['DB_SSL_MODE'])
    if args.db_name:
        gh_set(repo_full, 'DB_NAME', args.db_name)

    print("\nAll secrets uploaded to GitHub successfully.")


if __name__ == '__main__':
    main()
