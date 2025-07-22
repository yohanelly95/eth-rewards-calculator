#!/bin/bash

# Ethereum Rewards Calculator - Quick Start Script

echo "🚀 Setting up Ethereum Rewards Calculator..."

# Create project directory
PROJECT_NAME="eth-rewards-calculator"
mkdir -p $PROJECT_NAME
cd $PROJECT_NAME

# Initialize Go module
echo "📦 Initializing Go module..."
go mod init eth-rewards-calculator

# Create directory structure
echo "📁 Creating project structure..."
mkdir -p cmd/calculator
mkdir -p internal/calculator
mkdir -p internal/config
mkdir -p internal/types
mkdir -p bin

# Download dependencies
echo "⬇️  Installing dependencies..."
go get github.com/spf13/pflag@v1.0.5
go get github.com/fatih/color@v1.15.0

echo "✅ Project structure created!"
echo ""
echo "📋 Next steps:"
echo "1. Copy the provided Go files to their respective directories"
echo "2. Run 'make build' to compile the calculator"
echo "3. Run './bin/eth-rewards -validators 4096' to see an example"
echo ""
echo "🔧 Example commands:"
echo "  # Build the project"
echo "  make build"
echo ""
echo "  # Calculate for 4,096 validators"
echo "  ./bin/eth-rewards -validators 4096"
echo ""
echo "  # Compare different validator counts"
echo "  ./bin/eth-rewards -compare 1000,10000,100000,500000,1000000"
echo ""
echo "  # Show detailed breakdown with penalties"
echo "  ./bin/eth-rewards -validators 4096 -detailed -penalties"
echo ""
echo "📚 For more options, run: ./bin/eth-rewards --help"