#!/bin/bash

echo "mode: atomic" > coverage.out

if [ -f coverage_integration.out ]; then
    tail -n +2 coverage_integration.out >> coverage.out
    rm coverage_integration.out
fi

if [ -f coverage_unit.out ]; then
    tail -n +2 coverage_unit.out >> coverage.out
    rm coverage_unit.out
fi
