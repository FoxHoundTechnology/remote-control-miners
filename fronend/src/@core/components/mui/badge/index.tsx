import MuiBadge from '@mui/material/Badge'

import { CustomBadgeProps } from './types'

import useBgColor, { UseBgColorType } from 'src/@core/hooks/useBgColor'

const Badge = (props: CustomBadgeProps) => {
  const { sx, skin, color } = props

  const bgColors = useBgColor()

  const colors: UseBgColorType = {
    primary: { ...bgColors.primaryLight },
    secondary: { ...bgColors.secondaryLight },
    success: { ...bgColors.successLight },
    error: { ...bgColors.errorLight },
    warning: { ...bgColors.warningLight },
    info: { ...bgColors.infoLight }
  }

  return (
    <MuiBadge
      {...props}
      sx={skin === 'light' && color ? Object.assign({ '& .MuiBadge-badge': colors[color] }, sx) : sx}
    />
  )
}

export default Badge
