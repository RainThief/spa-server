#!/bin/bash
##############################################################################
#
# A script to run static analysis and all tests - will exit if any check fails
#
##############################################################################

# Assume this script is in the root directory
PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$PROJECT_DIR/scripts/include.sh"

_pushd "${PROJECT_DIR}"

set -e

echo_info "\nRunning Static Analysis"
./run_static_analysis.sh
echo_success "Static analysis passed\n"

echo_info "Running Unit Tests"
./run_unit_tests.sh
echo_success "Unit tests passed\n"

# system tests here

echo_success "All checks/tests successful"

_popd
