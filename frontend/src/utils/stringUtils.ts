/**
 * 字符串工具函数
 */

/**
 * 将驼峰命名法转换为蛇形命名法
 * 例如：chainType -> chain_type
 */
export function camelToSnakeCase(str: string): string {
  return str.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`);
}

/**
 * 将对象的键从驼峰命名法转换为蛇形命名法
 * 递归处理嵌套对象
 */
export function convertKeysToSnakeCase(obj: any): any {
  if (obj === null || typeof obj !== 'object') {
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.map(item => convertKeysToSnakeCase(item));
  }

  const result: Record<string, any> = {};
  
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      const snakeKey = camelToSnakeCase(key);
      const value = obj[key];
      
      // 特殊处理链类型
      if (key === 'chainType' && typeof value === 'string') {
        result[snakeKey] = value.toLowerCase();
      } else {
        result[snakeKey] = convertKeysToSnakeCase(value);
      }
    }
  }
  
  return result;
} 