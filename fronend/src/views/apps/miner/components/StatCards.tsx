import { Grid } from '@mui/material'
import StatsCard from './StatsCard'

import Icon from 'src/@core/components/icon'
import { MinerType } from 'src/types/apps/minerTypes'
import { useMemo } from 'react'

type StatsCardsProps = {
  filteredStore: MinerType[]
}

const StatsCards = ({ filteredStore }: StatsCardsProps) => {
  // IN PROGRESS
  // call the react use memo

  console.log('filteredStore: in CARD ', filteredStore)

  const fleetStats = useMemo(() => {
    const excludedStatuses = ['disabled', 'unracked']
    const totalWorkers = filteredStore.filter(miner => !excludedStatuses.includes(miner.status)).length
    const warning = filteredStore.filter(miner => miner.status.toLowerCase().includes('warning')).length
    const offline = filteredStore.filter(miner => miner.status.toLowerCase().includes('offline')).length
    return {
      totalHashrate: filteredStore.reduce(
        (acc, miner) => acc + Math.floor(Math.floor(miner.hashRate / 10) / 100) / 1000,
        0
      ),
      activeStatus: filteredStore.filter(miner => miner?.status === 'Okay').length,
      totalWorkers: totalWorkers,
      warning: warning,
      offline: offline
    }
  }, [filteredStore])

  return (
    <Grid item xs={12}>
      <Grid container spacing={6}>
        <Grid item xs={12} md={3} sm={6} key={1}>
          <StatsCard
            title='Total Hashrate'
            stats={String(Math.floor(fleetStats.totalHashrate * 100) / 100)}
            trendNumber='1243432'
            unit='Ph/s'
            trend='positive'
            icon={<Icon icon='mdi:chart-bell-curve-cumulative' />}
          />
        </Grid>
        <Grid item xs={12} md={3} sm={6} key={1}>
          <StatsCard
            title='Active Workers'
            stats={String(fleetStats.activeStatus)}
            trendNumber='1600'
            trend='positive'
            icon={<Icon icon='mdi:wifi-check' fontSize={20} />}
          />
        </Grid>
        <Grid item xs={12} md={3} sm={6} key={1}>
          <StatsCard
            title='Total Workers'
            stats={String(fleetStats.totalWorkers)}
            trendNumber='1243432'
            trend='positive'
            icon={<Icon icon='mdi:server' fontSize={20} />}
          />
        </Grid>
        <Grid item xs={12} md={3} sm={6} key={1}>
          <StatsCard
            title='Warning/Offline'
            stats={String(fleetStats.warning) + '/' + String(fleetStats.offline)}
            trendNumber='1243432'
            trend='positive'
            icon={<Icon icon='mdi:progress-alert' fontSize={20} />}
          />
        </Grid>
      </Grid>
    </Grid>
  )
}

export default StatsCards
