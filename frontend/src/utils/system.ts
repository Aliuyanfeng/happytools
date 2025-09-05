export interface SystemInfo {
  cpuUsage: number;
  memoryUsage: number;
  model: string;
  uploadSpeed: number;
  downloadSpeed: number;
}


let startTime: number;

// 初始化开始时间
const initializeStartTime = () => {
  startTime = Date.now();
};

// 获取系统信息（模拟实现）
export const getSystemInfo = async (): Promise<SystemInfo> => {
  // 在实际应用中，这些数据应该通过Wails调用Go后端获取
  // 这里使用随机数模拟系统信息
  return new Promise((resolve) => {
    setTimeout(async () => {
      resolve({
        cpuUsage: 12,
        memoryUsage: Math.floor(Math.random() * 100),
        model: 'Windows 10 Pro / Intel i7-10700K / 32GB RAM',
        uploadSpeed: parseFloat((Math.random() * 10).toFixed(2)),
        downloadSpeed: parseFloat((Math.random() * 50).toFixed(2))
      });
    }, 100);
  });
};

// 获取运行时间
export const getRuntime = (): string => {
  if (!startTime) {
    initializeStartTime();
  }

  const elapsedMs = Date.now() - startTime;
  const totalSeconds = Math.floor(elapsedMs / 1000);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  // 格式化时间为 HH:MM:SS
  return [
    hours.toString().padStart(2, '0'),
    minutes.toString().padStart(2, '0'),
    seconds.toString().padStart(2, '0')
  ].join(':');
};

// 初始化开始时间
initializeStartTime();