import * as moment from 'moment-timezone'
import dotenv from 'dotenv'
// Load environment variables from .env file
dotenv.config()

export function getAdjustedCurrentTime(): Date {
  console.log('getAdjustedCurrentTime TIMEZONE FROM ENV', process.env.TIME_ZONE)
  const timeZone = process.env.TIME_ZONE || 'UTC'
  const offsetHours = parseInt(process.env.OFFSET_HOURS || '0')

  // First, adjust time to the specified time zone.
  let adjustedTime = moment.tz(new Date(), timeZone)
  // Then, apply the additional offset.
  adjustedTime = adjustedTime.add(offsetHours, 'hours')

  return adjustedTime.toDate()
}

export const secondsToDHM = (seconds: number) => {
  const days = Math.floor(seconds / (3600 * 24))
  seconds -= days * 3600 * 24
  const hrs = Math.floor(seconds / 3600)
  seconds -= hrs * 3600
  const mnts = Math.floor(seconds / 60)

  return [days + 'd', hrs + 'h', mnts + 'm'].join(' ')
}

export const convertDateTime = (dateTimeStr: string) => {
  const date = new Date(dateTimeStr)
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0') // JS months are 0-indexed
  const day = date.getDate().toString().padStart(2, '0')
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')

  return {
    date: `${year}/${month}/${day}`,
    time: `${hour}:${minute}`
  }
}

export const convertToRegularDateAndTime = (timestamp: string) => {
  const dateObj = new Date(timestamp)

  // Extracting date details
  const year = dateObj.getFullYear()
  const month = String(dateObj.getMonth() + 1).padStart(2, '0') // Month index starts from 0
  const day = String(dateObj.getDate()).padStart(2, '0')

  // Extracting time details
  const hours = String(dateObj.getHours()).padStart(2, '0')
  const minutes = String(dateObj.getMinutes()).padStart(2, '0')
  const seconds = String(dateObj.getSeconds()).padStart(2, '0')

  // Setting up AM/PM format
  const amOrPm = dateObj.getHours() >= 12 ? 'PM' : 'AM'

  // Formatting the 24-hour time to 12-hour format
  const twelveHourFormat = Number(hours) % 12 || 12

  return `${twelveHourFormat}:${minutes}:${seconds} ${amOrPm} | ${month}/${day}/${year} `
}

// color gradient generator for the chart's legends
type RGB = [number, number, number]

const interpolateColor = (color1: RGB, color2: RGB, factor = 0.5): RGB => {
  const result: RGB = [...color1] as RGB
  for (let i = 0; i < 3; i++) {
    result[i] = Math.round(result[i] + factor * (color2[i] - color1[i]))
  }
  return result
}

export const getIntFromHex = (hex: string): RGB => {
  const bigint = parseInt(hex.slice(1), 16) // Remove the '#' at the start of the hex code
  return [(bigint >> 16) & 255, (bigint >> 8) & 255, bigint & 255]
}

export const getHexFromInt = (rgb: RGB): string => {
  return '#' + ((1 << 24) + (rgb[0] << 16) + (rgb[1] << 8) + rgb[2]).toString(16).slice(1).toUpperCase()
}

export const generateGradient = (startColor: string, endColor: string, steps: number): string[] => {
  const stepFactor = 1 / (steps - 1)
  const interpolatedColorArray: string[] = []
  for (let i = 0; i < steps; i++) {
    const interpolatedColor = interpolateColor(getIntFromHex(startColor), getIntFromHex(endColor), stepFactor * i)
    interpolatedColorArray.push(getHexFromInt(interpolatedColor))
  }
  return interpolatedColorArray
}

/* Usage example
  const lineColors = generateGradient('#FF5733', '#5733FF', 12)
*/
