#!/bin/bash
set -e

echo "🔄 Testing reproducible builds..."

# Create test directory
mkdir -p ./test-reproducible

# Build first version
echo "  Building first version..."
SOURCE_DATE_EPOCH=$(git log -1 --format=%ct) goreleaser release --snapshot --clean >/dev/null 2>&1 || {
    echo "❌ First build failed"
    goreleaser release --snapshot --clean
    rm -rf ./test-reproducible
    exit 1
}
cp -r ./dist ./test-reproducible/build1

# Wait to ensure different build times
sleep 2

# Build second version
echo "  Building second version..."
SOURCE_DATE_EPOCH=$(git log -1 --format=%ct) goreleaser release --snapshot --clean >/dev/null 2>&1 || {
    echo "❌ Second build failed"
    goreleaser release --snapshot --clean
    rm -rf ./test-reproducible
    exit 1
}
cp -r ./dist ./test-reproducible/build2

# Compare builds
echo "  Comparing builds..."
find ./test-reproducible/build1 -name "*.deb" -exec sha256sum {} \; | sed 's|./test-reproducible/build1/||' | sort > ./test-reproducible/checksums1.txt
find ./test-reproducible/build2 -name "*.deb" -exec sha256sum {} \; | sed 's|./test-reproducible/build2/||' | sort > ./test-reproducible/checksums2.txt

if diff ./test-reproducible/checksums1.txt ./test-reproducible/checksums2.txt >/dev/null 2>&1; then
    echo "✅ Builds are reproducible!"
else
    echo "❌ Builds are NOT reproducible!"
    echo "Differences:"
    diff ./test-reproducible/checksums1.txt ./test-reproducible/checksums2.txt || true
    rm -rf ./test-reproducible
    exit 1
fi

# Cleanup
rm -rf ./test-reproducible
echo "🎉 Reproducibility test passed"
