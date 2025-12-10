#!/bin/bash

# User-editable variables
VM_NAME="clab_cli"
ISO_PATH="$HOME/Downloads/ubuntu-24.04.3-live-server-arm64.iso"
VDI_PATH="$HOME/VirtualBox VMs/$VM_NAME/${VM_NAME}.vdi"
VDI_SIZE_MB=20480     # 20 GB

# Create the VM
VBoxManage createvm --name "${VM_NAME}" --ostype "Ubuntu24_LTS_arm64" --register

# Configure hardware
VBoxManage modifyvm "${VM_NAME}" --memory 16384 --cpus 16 --firmware efi --chipset armv8
VBoxManage modifyvm "${VM_NAME}" --graphicscontroller VMSVGA --vram 16

# Add VirtioSCSI storage controller
VBoxManage storagectl "${VM_NAME}" --name "VirtioSCSI" --add virtio-scsi --bootable on

# Create and attach virtual hard disk
VBoxManage createmedium disk --filename "${VDI_PATH}" --size "${VDI_SIZE_MB}"
VBoxManage storageattach "${VM_NAME}" --storagectl "VirtioSCSI" --port 0 --device 0 --type hdd --medium "${VDI_PATH}"

# Attach Ubuntu ARM64 ISO (for installation)
VBoxManage storageattach "${VM_NAME}" --storagectl "VirtioSCSI" --port 1 --device 0 --type dvddrive --medium "${ISO_PATH}"

# Enable basic NAT networking (you can change to bridged if you wish)
VBoxManage modifyvm "${VM_NAME}" --nic1 nat

# Start the VM
VBoxManage startvm "${VM_NAME}"

echo "VM '${VM_NAME}' created and started. Complete the OS install in the VirtualBox window."
echo "REMINDER: After installation, detach the ISO using:"
echo "  VBoxManage storageattach \"${VM_NAME}\" --storagectl \"VirtioSCSI\" --port 1 --device 0 --type dvddrive --medium none"