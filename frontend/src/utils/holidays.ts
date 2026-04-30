/**
 * 中国法定节假日工具
 * 数据来源：NateScarlet/holiday-cn (https://github.com/NateScarlet/holiday-cn)
 * 离线 JSON 文件，打包进前端，无需网络请求
 */

interface HolidayDay {
  name: string
  date: string      // YYYY-MM-DD
  isOffDay: boolean // true=放假, false=调休上班
}

interface HolidayYear {
  year: number
  days: HolidayDay[]
}

// 懒加载缓存：year -> Map<date, HolidayDay>
const cache = new Map<number, Map<string, HolidayDay>>()

// 动态 import 对应年份的 JSON
async function loadYear(year: number): Promise<Map<string, HolidayDay>> {
  if (cache.has(year)) return cache.get(year)!

  let data: HolidayYear | null = null
  try {
    // Vite 支持动态 import JSON
    const mod = await import(`../assets/holidays/${year}.json`)
    data = mod.default as HolidayYear
  } catch {
    // 该年份数据不存在，返回空 map
    const empty = new Map<string, HolidayDay>()
    cache.set(year, empty)
    return empty
  }

  const map = new Map<string, HolidayDay>()
  for (const day of data.days) {
    map.set(day.date, day)
  }
  cache.set(year, map)
  return map
}

export interface DayHolidayInfo {
  /** 是否是法定放假日（含周末调整后的假期） */
  isHoliday: boolean
  /** 是否是调休上班日（周末但要上班） */
  isWorkday: boolean
  /** 节日名称，如"春节"、"国庆节" */
  name: string
}

/**
 * 查询某一天的节假日信息
 * @param date YYYY-MM-DD
 */
export async function getDayHolidayInfo(date: string): Promise<DayHolidayInfo> {
  const year = parseInt(date.substring(0, 4))
  const map = await loadYear(year)
  const day = map.get(date)

  if (!day) {
    return { isHoliday: false, isWorkday: false, name: '' }
  }

  return {
    isHoliday: day.isOffDay,
    isWorkday: !day.isOffDay,
    name: day.name,
  }
}

/**
 * 批量查询一个月内所有日期的节假日信息
 * @param year  年份
 * @param month 月份 (1-12)
 * @returns Map<date, DayHolidayInfo>
 */
export async function getMonthHolidayMap(
  year: number,
  month: number
): Promise<Map<string, DayHolidayInfo>> {
  const map = await loadYear(year)
  const result = new Map<string, DayHolidayInfo>()

  // 只筛选该月的数据
  const prefix = `${year}-${String(month).padStart(2, '0')}`
  for (const [date, day] of map.entries()) {
    if (date.startsWith(prefix)) {
      result.set(date, {
        isHoliday: day.isOffDay,
        isWorkday: !day.isOffDay,
        name: day.name,
      })
    }
  }

  return result
}
