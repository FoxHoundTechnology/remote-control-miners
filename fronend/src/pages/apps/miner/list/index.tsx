import { useState, useEffect, MouseEvent } from 'react'

import Link from 'next/link'

import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import Menu from '@mui/material/Menu'
import Grid from '@mui/material/Grid'
import Divider from '@mui/material/Divider'
import { styled } from '@mui/material/styles'
import MenuItem from '@mui/material/MenuItem'
import IconButton from '@mui/material/IconButton'
import Typography from '@mui/material/Typography'
import Tooltip from '@mui/material/Tooltip'
import InputLabel from '@mui/material/InputLabel'
import FormControl from '@mui/material/FormControl'
import CardContent from '@mui/material/CardContent'
import { DataGrid, GridColDef, GridRowSelectionModel } from '@mui/x-data-grid'
import Select, { SelectChangeEvent } from '@mui/material/Select'

import Icon from 'src/@core/components/icon'

import CustomChip from 'src/@core/components/mui/chip'

import { ThemeColor } from 'src/@core/layouts/types'

import TableHeader from 'src/views/apps/miner/components/TableHeader'
import { ExtractFields, fetchMinerList } from 'src/store/apps/miner/list'

import { MinerType } from 'src/types/apps/minerTypes'

import { convertDateTime, secondsToDHM } from 'src/util'
import StatsCards from 'src/views/apps/miner/components/StatCards'
import { useMutation, useQuery } from 'react-query'
import { Command, sendCommand } from 'src/store/apps/minerControl'
import { Skeleton } from '@mui/material'

// TODO: fix the filter
// TODO: separte the skelton loader from the main table component
// TODO: stats card's slow rendering
// TODO: enum for status/label/color map
// TODO: seggregate the logic for remoteControlCallback into store management folder/component

// TODO: --  DUPLICATES --
interface ColorsType {
  [key: string]: ThemeColor
}

interface CellType {
  row: MinerType
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

// const minerModeObj: MinerModeType = {
//   '1': 'warning',
//   '3': 'info',
//   '0': 'success'
// }

const LinkStyled = styled(Link)(({ theme }) => ({
  fontWeight: 600,
  fontSize: '1rem',
  cursor: 'pointer',
  textDecoration: 'none',
  color: theme.palette.text.secondary,
  '&:hover': {
    color: theme.palette.primary.main
  }
}))

// NOTE:
// Row options will be
// normal, sleep, lowpower, reboot, disable, setting, details, pool
const RowOptions = ({ row }: any) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)

  const rowOptionsOpen = Boolean(anchorEl)

  const handleRowOptionsClick = (event: MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleRowOptionsClose = () => {
    setAnchorEl(null)
  }

  const handleReboot = () => {}
  const handleSleep = () => {}
  const handleUnrack = () => {}
  const handleDisable = () => {}
  const handleReactivate = () => {}
  const handleDelete = () => {}

  return (
    <>
      <IconButton size='small' onClick={handleRowOptionsClick}>
        <Icon icon='mdi:dots-vertical' />
      </IconButton>
      <Menu
        keepMounted
        anchorEl={anchorEl}
        open={rowOptionsOpen}
        onClose={handleRowOptionsClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right'
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'right'
        }}
        PaperProps={{ style: { minWidth: '8rem' } }}
      >
        <MenuItem
          component={Link}
          sx={{ '& svg': { mr: 2 } }}
          onClick={handleRowOptionsClose}
          href={`/apps/miner/${row?.macAddress}/view`}
        >
          <Icon icon='mdi:eye-outline' fontSize={20} />
          View
        </MenuItem>
        <MenuItem onClick={handleReboot} sx={{ '& svg': { mr: 2 } }}>
          <Icon icon='mdi:restart' fontSize={20} />
          Reboot
        </MenuItem>
        <MenuItem onClick={handleSleep} sx={{ '& svg': { mr: 2 } }}>
          <Icon icon='mdi:pause-circle-outline' fontSize={20} />
          Sleep
        </MenuItem>
        {row?.status === 'disabled' || row?.status === 'unracked' ? (
          <MenuItem onClick={handleReactivate} sx={{ '& svg': { mr: 2 } }}>
            <Icon icon='mdi:restart' fontSize={20} />
            Reactivate
          </MenuItem>
        ) : (
          <>
            <MenuItem disabled onClick={handleDisable} sx={{ '& svg': { mr: 2 } }}>
              <Icon icon='mdi:delete-outline' fontSize={20} />
              Disable
            </MenuItem>
            <MenuItem disabled onClick={handleDelete} sx={{ '& svg': { mr: 2 } }}>
              <Icon icon='mdi:wifi-remove' fontSize={20} />
              Unrack
            </MenuItem>
          </>
        )}
      </Menu>
    </>
  )
}

const columns: GridColDef[] = [
  {
    flex: 0.2,
    minWidth: 180,
    maxWidth: 180,
    field: 'macAddress',

    renderCell: ({ row }: CellType) => {
      const { model, macAddress, ip } = row

      // TODO: add the ip on hover
      return (
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <Tooltip title={`${ip}`} arrow>
            <Box sx={{ display: 'flex', alignItems: 'flex-start', flexDirection: 'column' }}>
              <LinkStyled href={`/apps/miner/${macAddress}/view`}>{model}</LinkStyled>
              <Typography noWrap variant='caption'>
                {`@${macAddress}`}
              </Typography>
            </Box>
          </Tooltip>
        </Box>
      )
    }
  },
  {
    flex: 0.2,
    minWidth: 180,
    field: 'client',
    headerName: 'Client / Location',
    renderCell: ({ row }: CellType) => {
      return (
        <Box>
          <Typography noWrap>{row?.client}</Typography>
          <Typography noWrap variant='caption'>
            {`${row?.location}`}
          </Typography>
        </Box>
      )
    }
  },
  {
    flex: 0.1,
    minWidth: 90,
    field: 'status',
    headerName: 'Status',
    renderCell: ({ row }: CellType) => {
      // TODO: fix the status
      const { status } = row
      return (
        <Tooltip title={''}>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <CustomChip
              skin='light'
              size='small'
              label={status}
              color={statusColors[row?.status]}
              sx={{ textTransform: 'capitalize' }}
            />
          </Box>
        </Tooltip>
      )
    }
  },
  {
    flex: 0.1,
    minWidth: 120,
    field: 'mode',
    headerName: 'Mode',
    renderCell: ({ row }: CellType) => {
      // online, offline, lowpower, sleep, reboot, disable, warning
      return (
        <CustomChip
          skin='light'
          size='small'
          label={row?.mode}
          color={modeColors[row?.mode]}
          sx={{ textTransform: 'capitalize' }}
        />
      )
    }
  },
  {
    flex: 0.1,
    minWidth: 140,
    field: 'hashRate',
    headerName: 'Hashrate',
    renderCell: ({ row }: CellType) => {
      const thFormat = row?.hashRate ? Math.floor(row?.hashRate / 10) / 100 : 0

      return (
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <Typography noWrap sx={{ color: 'text.primary', textTransform: 'capitalize' }}>
            {thFormat}{' '}
            <Typography ml={1} variant='caption'>
              TH/s
            </Typography>
          </Typography>
        </Box>
      )
    }
  },
  {
    flex: 0.1,
    minWidth: 80,
    field: 'maxTemp',
    headerName: 'Temp',
    renderCell: ({ row }: CellType) => {
      return (
        <Tooltip title={`${row?.tempArr}`} arrow>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Typography noWrap sx={{ color: 'text.secondary', textTransform: 'capitalize' }}>
              {row?.maxTemp}°C
            </Typography>
          </Box>
        </Tooltip>
      )
    }
  },
  {
    flex: 0.1,
    minWidth: 100,
    field: 'maxFan',
    headerName: 'Fan',
    renderCell: ({ row }: CellType) => {
      return (
        <Tooltip title={`${row.fanArr}`} arrow>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Typography noWrap sx={{ color: 'text.secondary', textTransform: 'capitalize' }}>
              {row?.maxFan}RPM
            </Typography>
          </Box>
        </Tooltip>
      )
    }
  },
  {
    flex: 0.1,
    minWidth: 120,
    field: 'upTime',
    headerName: 'Uptime',
    renderCell: ({ row }: CellType) => {
      const upTime = row?.upTime ? secondsToDHM(row.upTime) : 0

      return (
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <Typography noWrap sx={{ color: 'text.secondary', textTransform: 'capitalize' }}>
            {upTime}
          </Typography>
        </Box>
      )
    }
  },
  {
    // TODO: add the time difference setter
    flex: 0.2,
    minWidth: 170,
    field: 'lastUpdated',
    headerName: 'Last Updated',
    renderCell: ({ row }: CellType) => {
      const lastUpdatedDate = row?.lastUpdated ? convertDateTime(row.lastUpdated).date : 'N/A'
      const lastUpdatedTime = row?.lastUpdated ? convertDateTime(row.lastUpdated).time : 'N/A'
      return (
        <Typography noWrap variant='caption'>
          {`${lastUpdatedTime}`}{' '}
          <Typography noWrap variant='caption'>
            | {lastUpdatedDate}
          </Typography>
        </Typography>
      )
    }
  },
  {
    // config option
    flex: 0.1,
    minWidth: 110,
    sortable: false,
    field: 'actions',
    headerName: '',
    renderCell: ({ row }: CellType) => {
      return (
        <Box>
          <RowOptions row={row} />
        </Box>
      )
    }
  }
]

// TODO: fix the search bar
// TODO: separate queries from the components
const MinerList = () => {
  const { data, isLoading, isError } = useQuery('miners', async () => {
    const response = await fetchMinerList()
    return response
  })
  const store = data?.miners ?? []

  const initialFilters = {
    client: '',
    minerType: '',
    status: '',
    location: ''
  }

  const [filters, setFilters] = useState(initialFilters)
  const [value, setValue] = useState<string>('')

  // WIP:  add the modal config UI
  const [paginationModel, setPaginationModel] = useState({ page: 0, pageSize: 50 })
  const [rowSelectionModel, setRowSelectionModel] = useState<GridRowSelectionModel>([])
  const [openFilter, setOpenFilter] = useState<boolean>(false)

  const [filteredStore, setFilteredStore] = useState(store)
  // filter hander for the global search across all the columns
  const handleSearchFilter = (value: string) => {
    const lowercasedValue = value.toLowerCase().trim()
    if (lowercasedValue === '') {
      setFilteredStore(store)
    } else {
      const filteredData = store.filter((item: any) =>
        [
          // 'id',
          'macAddress',
          'ip',
          // 'name',
          // 'minerType',
          // 'serialNumber',
          'location',
          'firmware'
          // 'client',
          // 'fleetName',
          // 'ipRange'
        ].some((key: any) => {
          const val = item[key]
          return val !== undefined && val.toString().toLowerCase().includes(lowercasedValue)
        })
      )
      setFilteredStore(filteredData)
    }
  }

  const handleFilterChange = (filterKey: keyof typeof filters) => (event: SelectChangeEvent) => {
    const value = event.target.value

    //setFilters(prev => ({ ...prev, [filterKey]: value }))
    setFilters(prev => {
      // If user selected 'reset' or '', return initial value for that filter
      if (value === 'reset' || value === '') {
        return { ...prev, [filterKey]: initialFilters[filterKey] }
      }

      // Otherwise, update the filter normally
      else {
        return { ...prev, [filterKey]: value }
      }
    })
  }

  useEffect(() => {
    const lowercasedFilters = Object.keys(filters).reduce((acc: any, key: string) => {
      acc[key] = filters[key as keyof typeof filters].toLowerCase().trim()

      return acc
    }, {} as typeof filters)

    const filteredData = store?.filter((item: any) => {
      return Object.keys(lowercasedFilters).every((key: any) => {
        if (lowercasedFilters[key] === '' || lowercasedFilters[key] === 'reset') {
          return true
        }
        const itemValue = item[key]?.toString().toLowerCase().trim()

        return itemValue?.includes(lowercasedFilters[key])
      })
    })

    setFilteredStore(filteredData)

    // TODO: add the stats update logic goes here
  }, [filters])

  const handleResetFilter = () => {
    setFilters({
      client: '',
      minerType: '',
      status: '',
      location: ''
    })
  }

  const toggleFilter = () => {
    setOpenFilter(!openFilter)
    if (!openFilter) {
      handleResetFilter()
    }
  }
  /*
	  MacAddresses []string          `json:"mac_addresses"`
	  Mode         miner_domain.Mode `json:"mode"`
  */

  // NOTE: Controller related logics
  const controllerMutation = useMutation({
    mutationFn: sendCommand
  })

  // NOTE: idArray = rowSelectionModel
  const remoteControlCallback = async (idArray: number[], command: Command) => {
    const macAddressArray: string[] = []
    idArray.map((id: number) => {
      const selectedMiner = store.find(miner => miner.id === id)
      macAddressArray.push(selectedMiner.macAddress)
    })

    controllerMutation.mutate({
      MacAddressess: macAddressArray,
      Command: command
    })
  }

  return (
    <Grid container spacing={3}>
      <Grid item xs={12} padding={1}>
        <StatsCards filteredStore={filteredStore} />
      </Grid>
      <Grid item xs={12} padding={1}>
        <Card sx={{ mb: 6 }}>
          <CardContent>
            {/* First Row */}
            {openFilter && (
              <Grid container spacing={6}>
                <Grid item sm={3} xs={12}>
                  {/* NOTE: here goes the select component 
                          since the types of filters for miners will most likely remain the same
                          we will hard-code the filters here for now
                  */}
                  <FormControl fullWidth>
                    <InputLabel id='client-select'>Select Client</InputLabel>
                    <Select
                      fullWidth
                      value={filters.client}
                      id='select-client'
                      label='Select Client'
                      labelId='client-select'
                      onChange={handleFilterChange('client')}
                      inputProps={{ placeholder: 'Select Client' }}
                    >
                      <MenuItem value='reset'>Select Client</MenuItem>
                      {(ExtractFields(store)?.client ?? []).map(option => (
                        <MenuItem value={option} key={option}>
                          {option}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>
                <Grid item sm={3} xs={12}>
                  <FormControl fullWidth>
                    {/*
                        change the id here to Location
                    */}
                    <InputLabel id='location-select'>Select Location</InputLabel>
                    <Select
                      fullWidth
                      value={filters.location}
                      id='select-plan'
                      label='Select Plan'
                      labelId='plan-select'
                      onChange={handleFilterChange('location')}
                      inputProps={{ placeholder: 'Select Location' }}
                    >
                      <MenuItem value='reset'>Select Location</MenuItem>
                      {(ExtractFields(store)?.location ?? []).map(option => (
                        <MenuItem value={option} key={option}>
                          {option}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>
                <Grid item sm={3} xs={12}>
                  <FormControl fullWidth>
                    <InputLabel id='status-select'>Select Status</InputLabel>
                    <Select
                      fullWidth
                      value={filters.status}
                      id='select-status'
                      label='Select Status'
                      labelId='status-select'
                      onChange={handleFilterChange('status')}
                      inputProps={{ placeholder: 'Select Role' }}
                    >
                      <MenuItem value='reset'>Select Status</MenuItem>
                      {(ExtractFields(store)?.status ?? []).map(option => (
                        <MenuItem value={option} key={option}>
                          {option}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>
                <Grid item sm={3} xs={12}>
                  <FormControl fullWidth>
                    <InputLabel id='status-select'>Select Miner Type</InputLabel>
                    <Select
                      fullWidth
                      value={filters.minerType}
                      id='select-miner-type'
                      label='Select Miner Type'
                      labelId='miner-type-select'
                      onChange={handleFilterChange('minerType')}
                      inputProps={{ placeholder: 'Select Miner Type' }}
                    >
                      <MenuItem value='reset'>Select Miner Type</MenuItem>
                      {(ExtractFields(store)?.minerType ?? []).map(option => (
                        <MenuItem value={option} key={option}>
                          {option}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>
              </Grid>
            )}
          </CardContent>
          <Divider />
          <TableHeader
            toggleFilter={toggleFilter}
            openFilter={openFilter}
            setValue={setValue}
            value={value}
            handleFilter={handleSearchFilter}
            rowSelectionModel={rowSelectionModel}
            setRowSelectionModel={setRowSelectionModel}
            remoteControlCallback={remoteControlCallback}
          />
          {/*
              NOTE: first we will edit the data structure and  
          */}
          {isLoading && (
            <Grid container paddingX={6} paddingY={3}>
              <Grid item xs={12}>
                <Skeleton animation='wave' height={40} />
              </Grid>
              <Grid item xs={12}>
                <Skeleton animation='wave' height={40} />
              </Grid>
              <Grid item xs={12}>
                <Skeleton animation='wave' height={40} />
              </Grid>
              <Grid item xs={12}>
                <Skeleton animation='wave' height={40} />
              </Grid>
              <Grid item xs={12}>
                <Skeleton animation='wave' height={40} />
              </Grid>
            </Grid>
          )}
          {data && (
            <DataGrid
              autoHeight
              rows={store ?? []}
              columns={columns}
              disableRowSelectionOnClick
              checkboxSelection
              // NOTE: only select the id of each row
              onRowSelectionModelChange={newRowSelectionModel => {
                setRowSelectionModel(newRowSelectionModel)
              }}
              rowSelectionModel={rowSelectionModel}
              pageSizeOptions={[50, 100, 500]}
              paginationModel={paginationModel}
              onPaginationModelChange={setPaginationModel}
            />
          )}
        </Card>
      </Grid>
    </Grid>
  )
}

export const exportToCSV = (csvData: any, fileName: string) => {
  // *** WIP ***
}
export default MinerList
