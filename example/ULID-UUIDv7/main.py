import uuid6
from datetime import datetime, timezone


def generate_uuid_v7() -> str:
    """生成 UUID v7，去掉 - 连接符"""
    return str(uuid6.uuid7()).replace("-", "")


def parse_uuid_v7_time(uuid_str: str) -> datetime:
    """从 UUID v7 字符串反推生成时间（支持有无 - 连接符）"""
    clean = uuid_str.replace("-", "")
    if len(clean) != 32:
        raise ValueError("无效的 UUID 字符串")

    standard = f"{clean[:8]}-{clean[8:12]}-{clean[12:16]}-{clean[16:20]}-{clean[20:]}"
    u = uuid6.UUID(standard)

    # uuid6==2022.1.31 高48位存的是 unix_ms * 4096 / 1000
    # 反推毫秒：raw * 1000 / 4096
    raw = u.int >> 80
    unix_ms = raw * 1000 / 4096
    return datetime.fromtimestamp(unix_ms / 1000, tz=timezone.utc)


if __name__ == "__main__":
    # uuid = generate_uuid_v7()
    # print(f"UUID v7 : {uuid}")

    dt = parse_uuid_v7_time("069b671304d4788cad3a5daa0eb25cf3")
    print(f"生成时间(UTC) : {dt}")
    print(f"生成时间(本地): {dt.astimezone()}")
    partition_table_name = f"t_file_{dt.strftime('%Y%m%d')}"
    print(f"分区表名: {partition_table_name}")