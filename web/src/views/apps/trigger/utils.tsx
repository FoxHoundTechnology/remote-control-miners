import { Target, Action } from 'src/types/apps/triggerTypes'

export const generateValueLabel = (target: Target): string => {
  console.log('generateValueLabel ==>>', target)

  let label = ''

  if (target.type === 'hashrate') {
    label = `${target?.value}% of ideal hashrate`
  }

  if (target.type === 'temperature') {
    label = `${target?.value}°C`
  }

  if (target.type === 'fan_speed') {
    label = `${target?.value} RPM`
  }

  return label
}

export const generateTargetSentence = (target: Target): string => {
  let sentence = ''

  if (target.type === 'hashrate' && target.value && target.percentage) {
    sentence = `When ${target.percentage}% of the fleet has the ${target.value}% of ideal hashrate`
  }

  if (target.type === 'temperature' && target.value && target.percentage) {
    sentence = `When ${target.percentage}% of the fleet exceeds ${target.value}°C`
  }

  if (target.type === 'fan_speed' && target.value && target.percentage) {
    sentence = `When ${target.percentage}% of the fleet's fan speed exceeds ${target.value} RPM`
  }

  if (target.type === 'offline' && target.percentage) {
    sentence = `When ${target.percentage}% of the fleet goes offline`
  }

  if (target.type === 'missing_hashboard' && target.percentage) {
    sentence = `When ${target.percentage}% of the fleet has a missing hashboard`
  }

  return sentence
}

export const generateActionSentence = (action: Action): string => {
  let sentence = ''

  if (action.type === 'reboot') {
    sentence = `Reboot the miners`
  }

  if (action.type === 'normal_mode') {
    sentence = `Change to normal mode`
  }

  if (action.type === 'sleep_mode') {
    sentence = `Change to sleep mode`
  }

  // if (action.Type === 'change_pool' && action.value) {
  //   sentence = `Change to pool: ${action.value}`
  // }

  // Add similar conditions for other Action types if needed

  return sentence
}
