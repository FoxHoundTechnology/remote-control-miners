import { ReactNode } from 'react'

import { ThemeColor } from 'src/@core/layouts/types'
import { OptionsMenuType } from 'src/@core/components/option-menu/types'

export type CardStatsHorizontalProps = {
  title: string
  stats: string
  icon?: ReactNode
  unit?: string
  color?: ThemeColor
  trendNumber: string
  trend?: 'positive' | 'negative'
}

// NOTE: This is the card props for the
// miner table page
export type CardStatsVerticalProps = {
  title: string
  stats: string
  icon: ReactNode
  subtitle: string
  color?: ThemeColor
  trendNumber: string
  trend?: 'positive' | 'negative'
  optionsMenuProps?: OptionsMenuType
  // newly added
  unit?: string
}

export type CardStatsCharacterProps = {
  src: string
  title: string
  stats: string
  chipText: string
  trendNumber: string
  chipColor?: ThemeColor
  trend?: 'positive' | 'negative'
}
