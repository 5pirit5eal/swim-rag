#!/bin/bash
set -e

source .config.env
# //////////////////////////////////////////////////////////////////////////////
# START tasks

run() {
  uv run fastmcp run src/swim_rag_mcp/main.py \
    --transport http 
}

dev() {
  uv run fastmcp dev --server-spec src/swim_rag_mcp/main.py
}

docker-run() {
  local container_id=$1
  local port=${2:-"8080"}
  # local background=${3:-"-d"}
  docker run \
    -v ~/.config/gcloud/application_default_credentials.json:/gcp/creds.json \
    -p $port:8080 \
    -e GOOGLE_APPLICATION_CREDENTIALS=/gcp/creds.json \
    -e GOOGLE_CLOUD_PROJECT="$PROJECT_ID" \
    --env-file .config.env \
    $background \
    -i $container_id
}

install() {
  uv sync --group dev
}

format() {
  local files=${1:-'.'}
  uv run ruff format "$files"
  uv run ruff check --select I --fix "$files"
}

validate() {
  local files=${1:-'.'}
  local mypy_cache=${2:-'.mypy_cache'}
  uv run ruff check "$files"
  uv run ruff format "$files" --diff
  uv run ruff check --select I "$files"
  mkdir -p "$mypy_cache"
  echo "Running mypy..."
  uv run mypy "$files" --cache-dir "$mypy_cache"
  
  echo "Security check with bandit..."
  uv run bandit -c pyproject.toml -ll -x ./.venv/ -r "$files"
}

integration-test() { # ADAPT ACCORDING TO YOUR TEST REQUIREMENTS
  local port=${2:-8080}
  # Start emulator
  gcloud emulators firestore start --host-port="localhost:$port" > /dev/null 2>&1 &
  export FIRESTORE_EMULATOR_HOST="localhost:$port"
  sleep 10 # Wait for emulator to start
  uv run --group dev pytest tests/integration_tests/test_* "$*"
  result=$?
  # Stop emulator
  curl -d '' "localhost:$port/shutdown" > /dev/null 2>&1 &
  sleep 1
  return $result
}

unit-test() { # ADAPT ACCORDING TO YOUR TEST REQUIREMENTS
  uv run --group dev pytest tests/unit_tests/test_* "$*"
  result=$?
  sleep 1
  return $result
}

system-test() { # ADAPT ACCORDING TO YOUR TEST REQUIREMENTS
  uv run --group dev pytest tests/system_tests/test_* "$*"
  result=$?
  sleep 1
  return $result
}

create-identity-token() {
  gcloud auth print-identity-token
}

authenticate() {
  gcloud auth login --update-adc --no-launch-browser
}

activate() {
  # ADAPT IF WANTING TO USE E.G. CONDA
  gcloud config configurations activate "$PROJECT_ID"
  authenticate
  gcloud auth application-default set-quota-project "$PROJECT_ID"
  echo "SUCCESS: GOOGLE CLOUD CONFIGURATION ACTIVATED"
}

setup-gcloud() {
  local setup_wif=${1:-"false"}
  local google_account=${2:-""}
  echo "--- SETTING UP LOCAL GOOGLE CLOUD SDK CONFIGURATION ---"
  gcloud config configurations create "$PROJECT_ID"
  if [ -n "$google_account" ]; then
    gcloud config set account "$google_account"
  fi
  activate
  gcloud config set project "$PROJECT_ID"
  gcloud config set compute/region "$REGION"
  gcloud components install cloud-firestore-emulator

  if [ "$setup_wif" = "true" ]; then
    setup-wif
  fi
}




# END tasks
# //////////////////////////////////////////////////////////////////////////////
help() {
  echo "Usage: ./Taskfile.sh [task]"
  echo
  echo "Available tasks:"
  echo "  run                           Run the application locally."
  echo "  docker-run                    Run the application in a previously built Docker container."
  echo "  format                        Format the code using ruff."
  echo "  validate                      Perform code linting and formatting using ruff and mypy."
  echo "  install                       Install development dependencies."
  echo "  integration-test              Run integration tests using pytest."
  echo "  unit-test                     Run unit tests using pytest."
  echo "  system-test                   Run system tests using pytest."
  echo "  create-identity-token         Create an identity token for external request authentication."
  echo "  authenticate                  Authenticate to Google Cloud."
  echo "  activate                      Activate Google Cloud configuration."
  echo "  setup-gcloud                  Set up the Google Cloud settings (and Workload Identity Federation)."
  echo "  setup-wif                     Set up Workload Identity Federation."
  echo
  echo "If no task is provided, the default is to run the application."
}

# Check if the provided argument matches any of the functions
if [ -n "$1" ] && ! declare -f "$1" > /dev/null; then
  echo "Error: Unknown task '$1'"
  echo
  help  # Show help if the task is not recognized
  exit 1
fi

# Run application if no argument is provided
"${@:-run}"