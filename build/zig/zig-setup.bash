# use github api identify latest zig version. Yes, really.
VERSION=$(curl -s "https://api.github.com/repos/ziglang/zig/releases/latest" | jq -r '.tag_name')

TARGET=$(curl -s "https://ziglang.org/download/index.json" | jq --arg VERSION "$VERSION" '.[$VERSION]."x86_64-linux"')

TARGET_URL=$(echo $TARGET | jq -r ".tarball")
TARGET_HASH=$(echo $TARGET | jq -r ".shasum")

# I have trust issue
aria2c -x 8 -s 8 -j 8 --check-integrity=true --checksum=sha-256=$TARGET_HASH "$TARGET_URL"

FILENAME=zig-linux-x86_64-${VERSION}.tar.xz
tar -xf $FILENAME && mv zig-linux-x86_64-${VERSION} /zig && rm $FILENAME