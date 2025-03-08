#!/bin/bash
 
# 注意：此脚本需要使用 sudo 权限运行，例如: sudo ./install-service.sh

# 定义用户和用户组。
USER="nobody"
GROUP="nobody"
# 定义日志根目录前缀。
LOG_ROOT_PREFIX="/data/wwwroot/log.ppno.net/wwwroot/tool"

# 定义服务名称。
SERVICE_NAME="sms-bridge"

# 检查是否以 root 权限运行。
if [ "$EUID" -ne 0 ]; then
    echo "错误：此脚本需要使用 sudo 权限运行，例如: sudo $0"
    exit 1
fi

# 检查是否为 Linux 系统。
if [[ "$(uname)" != "Linux" ]]; then
    echo "错误：此脚本仅支持在 Linux 系统上运行。"
    exit 1
fi

# 检查是否已经安装过服务。
if systemctl list-units --full --all | grep -q "^$SERVICE_NAME.service"; then
    echo "服务已经安装过，无需重复安装。"
    exit 0
fi

# 获取当前脚本所在目录。
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 定义日志根目录。
LOG_ROOT_DIR="${LOG_ROOT_PREFIX}/${SERVICE_NAME}"
if [ ! -d "${LOG_ROOT_DIR}" ]; then
  mkdir -p "${LOG_ROOT_DIR}"
fi

# 修改日志目录权限。
chown -R ${USER}:${GROUP} "${LOG_ROOT_DIR}"
chmod -R 755 "${LOG_ROOT_DIR}"

# 安装服务。
cp "${SCRIPT_DIR}/${SERVICE_NAME}.service" /usr/lib/systemd/system/

# 启动服务。
systemctl enable ${SERVICE_NAME}

# 启动服务。
systemctl start ${SERVICE_NAME}

# 查看服务状态。
systemctl status ${SERVICE_NAME}