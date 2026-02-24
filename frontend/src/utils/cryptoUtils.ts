/**
 * 加密工具函数库
 * @author: LiuYanFeng
 * @date: 2026-02-12
 * @description: 提供base64、md5、sha256、sha512、sha1编解码以及utf-8编解码功能
 */

// ==================== Base64 编解码 ====================

/**
 * Base64 编码
 * @param input 输入字符串
 * @returns Base64 编码后的字符串
 */
export function base64Encode(input: string): string {
  try {
    // 对于 Unicode 字符串，先转换为 UTF-8
    const utf8Bytes = new TextEncoder().encode(input);
    let binary = '';
    for (let i = 0; i < utf8Bytes.length; i++) {
      binary += String.fromCharCode(utf8Bytes[i]);
    }
    return btoa(binary);
  } catch (error) {
    console.error('Base64 编码错误:', error);
    throw new Error('Base64 编码失败');
  }
}

/**
 * Base64 解码
 * @param input Base64 编码的字符串
 * @returns 解码后的字符串
 */
export function base64Decode(input: string): string {
  try {
    const binary = atob(input);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    return new TextDecoder().decode(bytes);
  } catch (error) {
    console.error('Base64 解码错误:', error);
    throw new Error('Base64 解码失败，请检查输入是否为有效的 Base64 编码');
  }
}

// ==================== 哈希函数 ====================

/**
 * 计算字符串的 MD5 哈希值
 * @param input 输入字符串
 * @returns MD5 哈希值（十六进制）
 */
export async function md5Hash(input: string): Promise<string> {
  try {
    const encoder = new TextEncoder();
    const data = encoder.encode(input);
    const hashBuffer = await crypto.subtle.digest('MD5', data);
    return bufferToHex(hashBuffer);
  } catch (error) {
    console.error('MD5 计算错误:', error);
    throw new Error('MD5 计算失败');
  }
}

/**
 * 计算字符串的 SHA-1 哈希值
 * @param input 输入字符串
 * @returns SHA-1 哈希值（十六进制）
 */
export async function sha1Hash(input: string): Promise<string> {
  try {
    const encoder = new TextEncoder();
    const data = encoder.encode(input);
    const hashBuffer = await crypto.subtle.digest('SHA-1', data);
    return bufferToHex(hashBuffer);
  } catch (error) {
    console.error('SHA-1 计算错误:', error);
    throw new Error('SHA-1 计算失败');
  }
}

/**
 * 计算字符串的 SHA-256 哈希值
 * @param input 输入字符串
 * @returns SHA-256 哈希值（十六进制）
 */
export async function sha256Hash(input: string): Promise<string> {
  try {
    const encoder = new TextEncoder();
    const data = encoder.encode(input);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    return bufferToHex(hashBuffer);
  } catch (error) {
    console.error('SHA-256 计算错误:', error);
    throw new Error('SHA-256 计算失败');
  }
}

/**
 * 计算字符串的 SHA-512 哈希值
 * @param input 输入字符串
 * @returns SHA-512 哈希值（十六进制）
 */
export async function sha512Hash(input: string): Promise<string> {
  try {
    const encoder = new TextEncoder();
    const data = encoder.encode(input);
    const hashBuffer = await crypto.subtle.digest('SHA-512', data);
    return bufferToHex(hashBuffer);
  } catch (error) {
    console.error('SHA-512 计算错误:', error);
    throw new Error('SHA-512 计算失败');
  }
}

// ==================== UTF-8 编解码 ====================

/**
 * UTF-8 编码（字符串转字节数组）
 * @param input 输入字符串
 * @returns UTF-8 编码的字节数组（十六进制表示）
 */
export function utf8Encode(input: string): string {
  try {
    const encoder = new TextEncoder();
    const bytes = encoder.encode(input);
    return bufferToHex(bytes);
  } catch (error) {
    console.error('UTF-8 编码错误:', error);
    throw new Error('UTF-8 编码失败');
  }
}

/**
 * UTF-8 解码（字节数组转字符串）
 * @param hexString 十六进制表示的字节数组
 * @returns 解码后的字符串
 */
export function utf8Decode(hexString: string): string {
  try {
    // 清理输入，移除空格和0x前缀
    const cleanHex = hexString.replace(/\s|0x/gi, '');
    
    // 验证十六进制格式
    if (!/^[0-9a-fA-F]*$/.test(cleanHex)) {
      throw new Error('输入不是有效的十六进制字符串');
    }
    
    // 确保长度为偶数
    const paddedHex = cleanHex.length % 2 === 0 ? cleanHex : '0' + cleanHex;
    
    // 转换为字节数组
    const bytes = new Uint8Array(paddedHex.length / 2);
    for (let i = 0; i < bytes.length; i++) {
      const byte = parseInt(paddedHex.substr(i * 2, 2), 16);
      if (isNaN(byte)) {
        throw new Error('无效的十六进制字节');
      }
      bytes[i] = byte;
    }
    
    const decoder = new TextDecoder('utf-8');
    return decoder.decode(bytes);
  } catch (error) {
    console.error('UTF-8 解码错误:', error);
    throw new Error('UTF-8 解码失败: ' + (error as Error).message);
  }
}

// ==================== 辅助函数 ====================

/**
 * 将 ArrayBuffer 或 Uint8Array 转换为十六进制字符串
 * @param buffer ArrayBuffer 或 Uint8Array
 * @returns 十六进制字符串
 */
function bufferToHex(buffer: ArrayBuffer | Uint8Array): string {
  const bytes = buffer instanceof Uint8Array ? buffer : new Uint8Array(buffer);
  const hexArray: string[] = [];
  
  for (let i = 0; i < bytes.length; i++) {
    const hex = bytes[i].toString(16).padStart(2, '0');
    hexArray.push(hex);
  }
  
  return hexArray.join('');
}

/**
 * 将十六进制字符串转换为 Uint8Array
 * @param hexString 十六进制字符串
 * @returns Uint8Array
 */
export function hexToBytes(hexString: string): Uint8Array {
  const cleanHex = hexString.replace(/\s|0x/gi, '');
  const paddedHex = cleanHex.length % 2 === 0 ? cleanHex : '0' + cleanHex;
  
  const bytes = new Uint8Array(paddedHex.length / 2);
  for (let i = 0; i < bytes.length; i++) {
    bytes[i] = parseInt(paddedHex.substr(i * 2, 2), 16);
  }
  
  return bytes;
}

/**
 * 格式化十六进制字符串，每2个字符加一个空格
 * @param hexString 十六进制字符串
 * @returns 格式化后的字符串
 */
export function formatHex(hexString: string): string {
  const cleanHex = hexString.replace(/\s/g, '');
  return cleanHex.match(/.{1,2}/g)?.join(' ') || '';
}

/**
 * 检查字符串是否为有效的 Base64
 * @param input 输入字符串
 * @returns 是否为有效的 Base64
 */
export function isValidBase64(input: string): boolean {
  try {
    // Base64 正则表达式
    const base64Regex = /^[A-Za-z0-9+/]*={0,2}$/;
    if (!base64Regex.test(input)) {
      return false;
    }
    
    // 尝试解码
    base64Decode(input);
    return true;
  } catch {
    return false;
  }
}

/**
 * 检查字符串是否为有效的十六进制
 * @param input 输入字符串
 * @returns 是否为有效的十六进制
 */
export function isValidHex(input: string): boolean {
  const cleanHex = input.replace(/\s|0x/gi, '');
  return /^[0-9a-fA-F]*$/.test(cleanHex);
}

/**
 * 复制文本到剪贴板
 * @param text 要复制的文本
 * @returns 是否复制成功
 */
export async function copyToClipboard(text: string): Promise<boolean> {
  try {
    await navigator.clipboard.writeText(text);
    return true;
  } catch (error) {
    console.error('复制到剪贴板失败:', error);
    
    // 备用方法
    try {
      const textArea = document.createElement('textarea');
      textArea.value = text;
      textArea.style.position = 'fixed';
      textArea.style.left = '-999999px';
      textArea.style.top = '-999999px';
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();
      const successful = document.execCommand('copy');
      document.body.removeChild(textArea);
      return successful;
    } catch (fallbackError) {
      console.error('备用复制方法也失败:', fallbackError);
      return false;
    }
  }
}

// ==================== 批量处理函数 ====================

/**
 * 批量计算所有哈希值
 * @param input 输入字符串
 * @returns 包含所有哈希值的对象
 */
export async function computeAllHashes(input: string): Promise<{
  md5: string;
  sha1: string;
  sha256: string;
  sha512: string;
}> {
  try {
    const [md5, sha1, sha256, sha512] = await Promise.all([
      md5Hash(input),
      sha1Hash(input),
      sha256Hash(input),
      sha512Hash(input),
    ]);
    
    return {
      md5,
      sha1,
      sha256,
      sha512,
    };
  } catch (error) {
    console.error('批量计算哈希值错误:', error);
    throw new Error('批量计算哈希值失败');
  }
}

/**
 * 批量进行所有编码
 * @param input 输入字符串
 * @returns 包含所有编码结果的对象
 */
export async function computeAllEncodings(input: string): Promise<{
  base64: string;
  utf8Hex: string;
  md5: string;
  sha1: string;
  sha256: string;
  sha512: string;
}> {
  try {
    const [base64, utf8Hex, md5, sha1, sha256, sha512] = await Promise.all([
      Promise.resolve(base64Encode(input)),
      Promise.resolve(utf8Encode(input)),
      md5Hash(input),
      sha1Hash(input),
      sha256Hash(input),
      sha512Hash(input),
    ]);
    
    return {
      base64,
      utf8Hex,
      md5,
      sha1,
      sha256,
      sha512,
    };
  } catch (error) {
    console.error('批量编码错误:', error);
    throw new Error('批量编码失败');
  }
}