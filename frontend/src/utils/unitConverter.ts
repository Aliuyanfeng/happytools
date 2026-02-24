/**
 * 单位转换工具函数
 * @author: LiuYanFeng
 * @date: 2026-02-12
 * @description: 提供字节转换、时间转换、时间差值计算等功能
 */

// ==================== 字节转换相关 ====================

/**
 * 字节单位定义
 */
export const BYTE_UNITS = [
  { value: 'B', label: '字节 (B)', factor: 1 },
  { value: 'KB', label: '千字节 (KB)', factor: 1024 },
  { value: 'MB', label: '兆字节 (MB)', factor: 1024 * 1024 },
  { value: 'GB', label: '吉字节 (GB)', factor: 1024 * 1024 * 1024 },
  { value: 'TB', label: '太字节 (TB)', factor: 1024 * 1024 * 1024 * 1024 },
  { value: 'PB', label: '拍字节 (PB)', factor: 1024 * 1024 * 1024 * 1024 * 1024 },
  { value: 'EB', label: '艾字节 (EB)', factor: 1024 * 1024 * 1024 * 1024 * 1024 * 1024 },
] as const;

export type ByteUnit = typeof BYTE_UNITS[number]['value'];

/**
 * 字节转换函数
 * @param value 数值
 * @param fromUnit 源单位
 * @param toUnit 目标单位
 * @returns 转换后的数值
 */
export function convertBytes(value: number, fromUnit: ByteUnit, toUnit: ByteUnit): number {
  if (value === 0) return 0;
  
  const fromUnitInfo = BYTE_UNITS.find(unit => unit.value === fromUnit);
  const toUnitInfo = BYTE_UNITS.find(unit => unit.value === toUnit);
  
  if (!fromUnitInfo || !toUnitInfo) {
    throw new Error('Invalid byte unit');
  }
  
  // 先将值转换为字节，再转换为目标单位
  const bytes = value * fromUnitInfo.factor;
  return bytes / toUnitInfo.factor;
}

/**
 * 格式化字节数为人类可读的字符串
 * @param bytes 字节数
 * @param decimals 小数位数，默认为2
 * @returns 格式化后的字符串
 */
export function formatBytes(bytes: number, decimals: number = 2): string {
  if (bytes === 0) return '0 B';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

// ==================== 时间转换相关 ====================

/**
 * 时间单位定义
 */
export const TIME_UNITS = [
  { value: 'ms', label: '毫秒 (ms)', factor: 1 },
  { value: 's', label: '秒 (s)', factor: 1000 },
  { value: 'min', label: '分钟 (min)', factor: 1000 * 60 },
  { value: 'h', label: '小时 (h)', factor: 1000 * 60 * 60 },
  { value: 'd', label: '天 (d)', factor: 1000 * 60 * 60 * 24 },
  { value: 'w', label: '周 (w)', factor: 1000 * 60 * 60 * 24 * 7 },
  { value: 'mo', label: '月 (mo)', factor: 1000 * 60 * 60 * 24 * 30 }, // 近似值
  { value: 'y', label: '年 (y)', factor: 1000 * 60 * 60 * 24 * 365 }, // 近似值
] as const;

export type TimeUnit = typeof TIME_UNITS[number]['value'];

/**
 * 时间转换函数
 * @param value 数值
 * @param fromUnit 源单位
 * @param toUnit 目标单位
 * @returns 转换后的数值
 */
export function convertTime(value: number, fromUnit: TimeUnit, toUnit: TimeUnit): number {
  if (value === 0) return 0;
  
  const fromUnitInfo = TIME_UNITS.find(unit => unit.value === fromUnit);
  const toUnitInfo = TIME_UNITS.find(unit => unit.value === toUnit);
  
  if (!fromUnitInfo || !toUnitInfo) {
    throw new Error('Invalid time unit');
  }
  
  // 先将值转换为毫秒，再转换为目标单位
  const milliseconds = value * fromUnitInfo.factor;
  return milliseconds / toUnitInfo.factor;
}

/**
 * 格式化时间为人类可读的字符串
 * @param milliseconds 毫秒数
 * @returns 格式化后的字符串
 */
export function formatTime(milliseconds: number): string {
  if (milliseconds === 0) return '0 毫秒';
  
  const seconds = Math.floor(milliseconds / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  const weeks = Math.floor(days / 7);
  const months = Math.floor(days / 30);
  const years = Math.floor(days / 365);
  
  if (years > 0) {
    return `${years} 年 ${Math.floor((days % 365) / 30)} 月 ${days % 30} 天`;
  } else if (months > 0) {
    return `${months} 月 ${days % 30} 天 ${hours % 24} 小时`;
  } else if (weeks > 0) {
    return `${weeks} 周 ${days % 7} 天 ${hours % 24} 小时`;
  } else if (days > 0) {
    return `${days} 天 ${hours % 24} 小时 ${minutes % 60} 分钟`;
  } else if (hours > 0) {
    return `${hours} 小时 ${minutes % 60} 分钟 ${seconds % 60} 秒`;
  } else if (minutes > 0) {
    return `${minutes} 分钟 ${seconds % 60} 秒`;
  } else if (seconds > 0) {
    return `${seconds} 秒 ${milliseconds % 1000} 毫秒`;
  } else {
    return `${milliseconds} 毫秒`;
  }
}

// ==================== 时间差值计算相关 ====================

/**
 * 时间差值计算结果
 */
export interface TimeDifferenceResult {
  totalMilliseconds: number;
  formatted: string;
  breakdown: {
    years: number;
    months: number;
    days: number;
    hours: number;
    minutes: number;
    seconds: number;
    milliseconds: number;
  };
}

/**
 * 计算两个日期之间的时间差
 * @param startDate 开始日期
 * @param endDate 结束日期
 * @returns 时间差值结果
 */
export function calculateTimeDifference(startDate: Date, endDate: Date): TimeDifferenceResult {
  // 确保开始日期早于结束日期
  const start = startDate < endDate ? startDate : endDate;
  const end = startDate < endDate ? endDate : startDate;
  
  const totalMilliseconds = end.getTime() - start.getTime();
  
  // 计算各个时间单位
  const totalSeconds = Math.floor(totalMilliseconds / 1000);
  const totalMinutes = Math.floor(totalSeconds / 60);
  const totalHours = Math.floor(totalMinutes / 60);
  const totalDays = Math.floor(totalHours / 24);
  
  const years = Math.floor(totalDays / 365);
  const months = Math.floor((totalDays % 365) / 30);
  const days = totalDays % 30;
  const hours = totalHours % 24;
  const minutes = totalMinutes % 60;
  const seconds = totalSeconds % 60;
  const milliseconds = totalMilliseconds % 1000;
  
  // 构建格式化字符串
  let formatted = '';
  if (years > 0) formatted += `${years}年`;
  if (months > 0) formatted += `${months}月`;
  if (days > 0) formatted += `${days}天`;
  if (hours > 0) formatted += `${hours}小时`;
  if (minutes > 0) formatted += `${minutes}分钟`;
  if (seconds > 0) formatted += `${seconds}秒`;
  if (milliseconds > 0 && formatted === '') formatted += `${milliseconds}毫秒`;
  
  if (formatted === '') formatted = '0秒';
  
  return {
    totalMilliseconds,
    formatted,
    breakdown: {
      years,
      months,
      days,
      hours,
      minutes,
      seconds,
      milliseconds,
    },
  };
}

/**
 * 解析时间字符串为Date对象
 * @param timeString 时间字符串，支持格式：YYYY-MM-DD HH:mm:ss 或 HH:mm:ss
 * @returns Date对象
 */
export function parseTimeString(timeString: string): Date {
  const now = new Date();
  
  // 尝试解析完整日期时间格式
  const fullDateMatch = timeString.match(/^(\d{4})-(\d{2})-(\d{2})[T\s](\d{2}):(\d{2}):(\d{2})$/);
  if (fullDateMatch) {
    const [, year, month, day, hour, minute, second] = fullDateMatch;
    return new Date(
      parseInt(year),
      parseInt(month) - 1,
      parseInt(day),
      parseInt(hour),
      parseInt(minute),
      parseInt(second)
    );
  }
  
  // 尝试解析日期格式
  const dateMatch = timeString.match(/^(\d{4})-(\d{2})-(\d{2})$/);
  if (dateMatch) {
    const [, year, month, day] = dateMatch;
    return new Date(parseInt(year), parseInt(month) - 1, parseInt(day));
  }
  
  // 尝试解析时间格式
  const timeMatch = timeString.match(/^(\d{2}):(\d{2}):(\d{2})$/);
  if (timeMatch) {
    const [, hour, minute, second] = timeMatch;
    const date = new Date();
    date.setHours(parseInt(hour), parseInt(minute), parseInt(second), 0);
    return date;
  }
  
  // 尝试解析简单时间格式
  const simpleTimeMatch = timeString.match(/^(\d{2}):(\d{2})$/);
  if (simpleTimeMatch) {
    const [, hour, minute] = simpleTimeMatch;
    const date = new Date();
    date.setHours(parseInt(hour), parseInt(minute), 0, 0);
    return date;
  }
  
  // 如果都无法解析，返回当前时间
  return now;
}

/**
 * 格式化日期时间为字符串
 * @param date Date对象
 * @param format 格式字符串，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns 格式化后的字符串
 */
export function formatDateTime(date: Date, format: string = 'YYYY-MM-DD HH:mm:ss'): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  
  return format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds);
}