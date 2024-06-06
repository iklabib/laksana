# TODO: make apparmor and selinux profile to fix remount issue
docker run --rm -it -e BASE_URL=:8080 -p 8080:8080 \
           --cap-add=sys_admin \
           --cap-add=sys_chroot \
           --cap-add=sys_resource \
           --security-opt apparmor=unconfined \
           --security-opt label=disable quay.io/iklabib/markisa