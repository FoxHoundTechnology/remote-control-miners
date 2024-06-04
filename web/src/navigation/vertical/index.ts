/* eslint-disable lines-around-comment */

import { VerticalNavItemsType } from 'src/@core/layouts/types'

const navigation = (): VerticalNavItemsType => {
  return [
    {
      sectionTitle: 'Remote Control'
    },
    {
      title: 'Miners',
      icon: 'mdi:server',
      path: '/apps/miner/list'
    },
    {
      title: 'Scanner',
      icon: 'mdi:monitor-share'
      // path: '/apps/scanner/list'
    },
    {
      title: 'Miner Registration',
      icon: 'mdi:server-plus-outline'
      // path: '/apps/scanner'
    },
    {
      icon: 'mdi:bell-plus-outline',
      title: 'Trigger'
      // path: '/apps/trigger'
    },
    {
      title: 'Pools',
      icon: 'mdi:waves'
    },

    {
      sectionTitle: 'Power Control'
    },
    {
      title: 'Power',
      icon: 'mdi:transmission-tower'
    },
    {
      sectionTitle: 'Financial'
    },
    {
      title: 'Earnings',
      icon: 'mdi:wallet'
    },
    {
      sectionTitle: 'Tickets'
    },
    {
      title: 'Tickets',
      icon: 'mdi:ticket-confirmation-outline'
    },
    {
      sectionTitle: 'Settings'
    },
    {
      title: 'User Settings',
      icon: 'mdi:cog'
    }
  ]
}

export default navigation
