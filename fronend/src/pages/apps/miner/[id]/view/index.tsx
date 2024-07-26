import { useQuery } from 'react-query'
import { SyntheticEvent, useMemo, useState } from 'react'

import { Box, Card, CardContent, Divider, Link, Tab, Typography } from '@mui/material'
import LaunchIcon from '@mui/icons-material/Launch'
import MuiTabList, { TabListProps } from '@mui/lab/TabList'
import { MinerType } from 'src/types/apps/minerTypes'
import { styled, useTheme } from '@mui/material/styles'

import CustomChip from 'src/@core/components/mui/chip'
import { ThemeColor } from 'src/@core/layouts/types'
import { convertToRegularDateAndTime } from 'src/util'
import { MinerTimeSeriesDataResponse, PoolTimeSeriesDataResponse } from 'src/types/apps/minerDetailsTypes'
import TabContext from '@mui/lab/TabContext'
import Grid from '@mui/material/Grid'
import CircularProgress from '@mui/material/CircularProgress'

// icons
import ShowChartIcon from '@mui/icons-material/ShowChart'
import DeviceThermostatIcon from '@mui/icons-material/DeviceThermostat'
import SpeedIcon from '@mui/icons-material/Speed'
import WaterIcon from '@mui/icons-material/Water'
import NotificationAddIcon from '@mui/icons-material/NotificationAdd'
import SettingsIcon from '@mui/icons-material/Settings'

import HashrateChart from 'src/views/charts/hashrateChart'
import TemperatureChart from 'src/views/charts/temperatureChart'
import FanSpeedChart from 'src/views/charts/fanSpeedChart'
import PoolStatsChart from 'src/views/charts/poolStatsChart'
import NotificationSettings from 'src/views/apps/miner/components/NotificationSettings'

import RecentPoolStatsChart from 'src/views/charts/recentPoolStatsChart'
import { fetchMinerInfo, fetchMinerStats, fetchPoolStats } from 'src/store/apps/minerDetails'
import TabPanel from '@mui/lab/TabPanel'
import ActivityLogs from 'src/views/apps/miner/components/activityLogs'

// TODO: R&D for server component
// TODO: R&D for periodic requests with different params, using react-query and zustand
// TODO: Separation of concerns with view/layout/template
// TODO: Skelton loading screen
// TODO: seggregate tab logic to component folder

// TODO: --  DUPLICATES --
interface ColorsType {
  [key: string]: ThemeColor
}

const modeColors: ColorsType = {
  Normal: 'success',
  Sleep: 'warning',
  'Low Power': 'info',
  'N/A': 'warning'
}

const statusColors: ColorsType = {
  Online: 'primary',
  Offline: 'error',
  Disabled: 'warning',
  'Hashrate Error': 'error',
  'Temperature Error': 'error',
  'FanSpeed Error': 'error',
  'Missing Hashboard Error': 'error',
  'PoolShare Error': 'error'
}

// TODO: -- end of DUPLICATE --
// TODO: param for interval hours with react-query
type minerDetailsViewProps = {
  macAddress: string
}

// TODO: R&D for layout/template
// NOTE: logics related to   to right-half of the view
const TabList = styled(MuiTabList)<TabListProps>(({ theme }) => ({
  '& .MuiTabs-indicator': {
    display: 'none'
  },
  '& .Mui-selected': {
    backgroundColor: theme.palette.primary.main,
    color: `${theme.palette.common.white} !important`
  },
  '& .MuiTab-root': {
    minWidth: 65,
    minHeight: 40,
    paddingTop: theme.spacing(2),
    paddingBottom: theme.spacing(2),
    borderRadius: theme.shape.borderRadius,
    [theme.breakpoints.up('md')]: {
      minWidth: 130
    }
  }
}))

const minerDetailsView = ({ macAddress }: minerDetailsViewProps) => {
  const theme = useTheme()

  // NOTE; Tab related logics
  const [activeTab, setActiveTab] = useState<string>('hashrateChart')

  const handleChange = (event: SyntheticEvent, value: string) => {
    setActiveTab(value)
  }

  // Use React Query to fetch data
  const minerDetailsQuery = useQuery('minerDetails', () => fetchMinerInfo(macAddress), {
    staleTime: 0,
    cacheTime: 0
  })

  console.log('miner detail data in view: ', minerDetailsQuery.data)

  const minerStatsQuery = useQuery<MinerTimeSeriesDataResponse>(
    'minerStats',
    () => fetchMinerStats(macAddress, 24, 'h', 1, 'h'),
    {
      staleTime: 0,
      cacheTime: 0
    }
  )

  const poolStatsQuery = useQuery<PoolTimeSeriesDataResponse>(
    'poolStats',
    () => fetchPoolStats(macAddress, 24, 'h', 1, 'h'),
    {
      staleTime: 0,
      cacheTime: 0
    }
  )

  // const [isLoading, setIsLoading] = useState<boolean>(false)

  // Aggregate loading and error states
  const isLoading = minerDetailsQuery.isLoading || minerStatsQuery.isLoading || poolStatsQuery.isLoading
  const isError =
    minerDetailsQuery.isError ||
    minerStatsQuery.isError ||
    poolStatsQuery.isError ||
    !minerDetailsQuery.data ||
    !minerStatsQuery.data ||
    !poolStatsQuery.data

  // Process data only when all queries have succeeded
  const { hashrateArr, tempSensorArr, fanSensorArr, timestampArr } = useMemo(() => {
    if (isLoading || isError) {
      return { hashrateArr: [], tempSensorArr: [], fanSensorArr: [], timestampArr: [] }
    }

    const hashrateArr: number[] = []
    const tempSensorArr: number[][] = []
    const fanSensorArr: number[][] = []
    const timestampArr: string[] = []

    // First, ensure that miner_time_series_record exists
    if (minerStatsQuery?.data?.miner_time_series_record) {
      for (const record of minerStatsQuery?.data.miner_time_series_record) {
        hashrateArr.push(record.hashrate)
        tempSensorArr.push(record.temp_sensor)
        fanSensorArr.push(record.fan_sensor)
      }
    }

    if (minerStatsQuery?.data?.timestamps) {
      timestampArr.push(...minerStatsQuery.data.timestamps.map(String))
    }

    return {
      hashrateArr,
      tempSensorArr,
      fanSensorArr,
      timestampArr
    }
  }, [minerStatsQuery?.data, isLoading, isError])

  if (isLoading) {
    return <div>Loading...</div>
  }

  if (isError) {
    return <div>Error loading data. Please try again later.</div>
  }

  // useEffect(() => {}, [activeTab])

  return (
    <Grid container spacing={6}>
      <Grid item xs={12} md={5} lg={4}>
        <Grid container spacing={6}>
          <Grid item xs={12}>
            <Card>
              <CardContent sx={{ pt: 10, display: 'flex', alignItems: 'center', flexDirection: 'column' }}>
                {theme.palette.mode === 'dark' ? (
                  <Logo src={`/images/logos/antminer_white.png`} alt={''} />
                ) : (
                  <Logo src={`/images/logos/antminer_black.png`} alt={''} />
                )}
                <Typography variant='h5' sx={{ mt: 1, mb: 4 }}>
                  {minerDetailsQuery?.data?.model || 'N/A'}
                </Typography>
                <Box
                  sx={{
                    width: '100%',
                    display: 'flex',
                    flexDirection: 'row', // Changed from 'column' to 'row'
                    alignItems: 'center',
                    justifyContent: 'center', // This will align items to the start of the box
                    gap: 2 // This will add spacing between the two boxes
                  }}
                >
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <Typography sx={{ fontSize: '13px', marginRight: 2 }}>Status:</Typography>
                    <CustomChip
                      skin='light'
                      size='small'
                      label={minerDetailsQuery?.data?.status}
                      color={statusColors[minerDetailsQuery?.data?.status]}
                      sx={{ textTransform: 'capitalize' }}
                    />
                  </Box>
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <Typography sx={{ fontSize: '13px', marginRight: 2 }}>Mode:</Typography>
                    <CustomChip
                      skin='light'
                      size='small'
                      label={minerDetailsQuery?.data?.mode}
                      color={modeColors[minerDetailsQuery?.data?.mode]}
                      sx={{ textTransform: 'capitalize' }}
                    />
                  </Box>
                </Box>
              </CardContent>
              <CardContent>
                <Divider sx={{ mb: 4 }} />
                <Box sx={{ pb: 1 }}>
                  <Box sx={{ display: 'flex', mb: 2 }}>
                    <Typography sx={{ mr: 2, fontWeight: 500, fontSize: '0.875rem' }}>Client:</Typography>
                    <Typography variant='body2'>{minerDetailsQuery?.data?.client}</Typography>
                  </Box>
                  <Box sx={{ display: 'flex', mb: 2 }}>
                    <Typography sx={{ mr: 2, fontWeight: 500, fontSize: '0.875rem' }}>IP Address:</Typography>
                    <Link
                      href={`http://${minerDetailsQuery?.data?.ip}`}
                      underline='none'
                      sx={{ display: 'flex', alignItems: 'center' }}
                    >
                      <Typography variant='body2' sx={{ fontWeight: 600 }}>
                        {minerDetailsQuery?.data?.ip}
                      </Typography>
                      <LaunchIcon fontSize='small' sx={{ ml: 1 }} />
                    </Link>
                  </Box>
                  <Box sx={{ display: 'flex', mb: 2 }}>
                    <Typography sx={{ mr: 2, fontWeight: 500, fontSize: '0.875rem' }}>Mac Address:</Typography>
                    <Typography variant='body2'>{minerDetailsQuery?.data?.macAddress}</Typography>
                  </Box>
                  <Box sx={{ display: 'flex', mb: 2 }}>
                    <Typography sx={{ mr: 2, fontWeight: 500, fontSize: '0.875rem' }}>Serial Number:</Typography>
                    <Typography variant='body2' sx={{ textTransform: 'capitalize' }}>
                      {minerDetailsQuery?.data?.serialNumber}
                    </Typography>
                  </Box>
                  <Box sx={{ display: 'flex', mb: 2 }}>
                    <Typography sx={{ mr: 2, fontWeight: 500, fontSize: '0.875rem' }}>Location:</Typography>
                    <Typography variant='body2' sx={{ textTransform: 'capitalize' }}>
                      {minerDetailsQuery?.data?.location}
                    </Typography>
                  </Box>
                  <Box sx={{ display: 'flex', mb: 2 }}>
                    <Typography sx={{ mr: 2, fontWeight: 500, fontSize: '0.875rem' }}>Last Updated:</Typography>
                    <Typography variant='body2'>
                      {convertToRegularDateAndTime(minerDetailsQuery?.data?.lastUpdated)}
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12}>
            {!isLoading && !isError && poolStatsQuery?.data && <RecentPoolStatsChart poolStats={poolStatsQuery.data} />}
          </Grid>
        </Grid>{' '}
      </Grid>
      {/* NOTE: right half */}
      <Grid item xs={12} md={7} lg={8}>
        <Box>
          {!isLoading && !isError ? (
            <TabContext value={activeTab}>
              <TabList
                variant='scrollable'
                scrollButtons='auto'
                onChange={handleChange}
                aria-label='forced scroll tabs example'
              >
                <Tab
                  value='hashrateChart'
                  label={
                    <Box
                      sx={{
                        display: 'flex',
                        alignItems: 'center',
                        '& svg': { mr: 2 }
                      }}
                    >
                      <ShowChartIcon />
                      Hashrate
                    </Box>
                  }
                />
                <Tab
                  value='temperature'
                  label={
                    <Box sx={{ display: 'flex', alignItems: 'center', '& svg': { mr: 2 } }}>
                      <DeviceThermostatIcon />
                      Temperature
                    </Box>
                  }
                />
                <Tab
                  value='fanSpeed'
                  label={
                    <Box sx={{ display: 'flex', alignItems: 'center', '& svg': { mr: 2 } }}>
                      <SpeedIcon />
                      Fan Speed
                    </Box>
                  }
                />
                <Tab
                  value='pool'
                  label={
                    <Box sx={{ display: 'flex', alignItems: 'center', '& svg': { mr: 2 } }}>
                      <WaterIcon />
                      Pool
                    </Box>
                  }
                />
                <Tab
                  value='notification'
                  label={
                    <Box sx={{ display: 'flex', alignItems: 'center', '& svg': { mr: 2 } }}>
                      <NotificationAddIcon />
                      Notification
                    </Box>
                  }
                />
              </TabList>
              <Grid container sx={{ mt: 4 }}>
                <Grid item xs={12}>
                  <TabPanel value='hashrateChart'>
                    <HashrateChart hashrateArr={hashrateArr} timestampArr={timestampArr} />
                  </TabPanel>
                  <TabPanel value='temperature'>
                    <TemperatureChart tempSensorArr={tempSensorArr} timeStampArr={timestampArr} />
                  </TabPanel>
                  <TabPanel value='fanSpeed'>
                    <FanSpeedChart fanSensorArr={fanSensorArr} timeStampArr={timestampArr} />
                  </TabPanel>
                  <TabPanel value='pool'>
                    <PoolStatsChart
                      poolStatsArr={poolStatsQuery?.data?.pool_time_series_record}
                      timeStampArr={timestampArr}
                    />
                  </TabPanel>
                  <TabPanel value='notification'>
                    <NotificationSettings />
                  </TabPanel>
                </Grid>
                <Grid item xs={12} mt={6}>
                  <ActivityLogs />
                </Grid>
              </Grid>
            </TabContext>
          ) : isLoading ? (
            <div>Loading...</div>
          ) : (
            <div>Error loading data. Please try again later.</div>
          )}
        </Box>
      </Grid>
    </Grid>
  )
}

export async function getServerSideProps(context: any) {
  const { id } = context.query

  return {
    props: {
      macAddress: id
    }
  }
}

const Logo = styled('img')({
  width: '60%',
  height: 'auto',
  right: 30,
  left: 30
})

export default minerDetailsView
