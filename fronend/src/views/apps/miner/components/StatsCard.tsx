import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import { styled, useTheme } from '@mui/material/styles'
import Typography from '@mui/material/Typography'
import CardContent from '@mui/material/CardContent'
import MuiAvatar, { AvatarProps } from '@mui/material/Avatar'
import { alpha } from '@mui/material/styles'

import Icon from 'src/@core/components/icon'

import { CardStatsHorizontalProps } from 'src/@core/components/card-statistics/types'
import { ReactNode } from 'react'
import { ThemeColor } from 'src/@core/layouts/types'

const Avatar = styled(MuiAvatar)<AvatarProps>(({ theme }) => ({
  width: 44,
  height: 44,
  boxShadow: 'none',
  marginRight: theme.spacing(2.75),
  backgroundColor: theme.palette.background.paper,
  border: `0.5px solid ${theme.palette.primary.main}`,
  '& svg': {
    fontSize: '1.75rem'
  }
}))

export type StatsCardProps = {
  title: string
  stats: string
  icon?: ReactNode
  unit?: string
  color?: ThemeColor
  trendNumber: string
  trend?: 'positive' | 'negative'
}

const StatsCard = (props: CardStatsHorizontalProps) => {
  const { title, icon, stats, trendNumber, color = 'primary', trend = 'positive', unit } = props

  const theme = useTheme()

  return (
    <Card
      sx={{
        //backgroundColor: 'transparent !important',
        boxShadow: theme => `${theme.shadows[0]} !important`,
        border: `0.5px solid ${alpha(theme.palette.primary.main, 0.6)}`,
        paddingTop: '0px'
      }}
    >
      <CardContent>
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <Avatar
            variant='rounded'
            sx={{
              border: `1px solid ${alpha(theme.palette.primary.main, 0.95)}`,
              color: `${color}.main`,
              backgroundColor: `${color}.lighter`
            }}
          >
            {icon}
          </Avatar>
          <Box sx={{ display: 'flex', flexDirection: 'column' }}>
            <Typography variant='caption'>{title}</Typography>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', alignItems: 'center' }}>
              <Typography
                variant='h6'
                sx={{
                  mr: 1,
                  fontWeight: 600,
                  lineHeight: 1.05,
                  color: theme.palette.mode === 'light' ? theme.palette.grey[900] : theme.palette.primary.contrastText
                }}
              >
                {stats}
              </Typography>
              <Box
                sx={{
                  display: 'flex',
                  alignItems: 'center'
                }}
              >
                <Typography
                  variant='h6'
                  sx={{
                    mr: 1,
                    fontWeight: 600,
                    lineHeight: 1.05,
                    color: theme.palette.mode === 'light' ? theme.palette.grey[900] : theme.palette.primary.contrastText
                  }}
                >
                  {/* NOTE: unit number goes here */}
                  {unit}
                </Typography>
                {/* 
                <Box
                  component='span'
                  sx={{ display: 'inline-flex', color: trend === 'positive' ? 'success.main' : 'error.main' }}
                >
                  <Icon icon={trend === 'positive' ? 'mdi:chevron-up' : 'mdi:chevron-down'} />
                </Box>
                <Typography
                  variant='caption'
                  sx={{ fontWeight: 600, color: trend === 'positive' ? 'success.main' : 'error.main' }}
                >
                  {trendNumber}
                </Typography 
                */}
              </Box>
            </Box>
          </Box>
        </Box>
      </CardContent>
    </Card>
  )
}

export default StatsCard
