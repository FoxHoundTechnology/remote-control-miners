import { useCallback, useState, MouseEvent } from 'react'

import Grid from '@mui/material/Grid'
import IconButton from '@mui/material/IconButton'
import Typography from '@mui/material/Typography'
import { GridColDef } from '@mui/x-data-grid'
import { SelectChangeEvent } from '@mui/material/Select'

import Icon from 'src/@core/components/icon'

import CustomChip from 'src/@core/components/mui/chip'

import { ThemeColor } from 'src/@core/layouts/types'

import TableHeader from 'src/views/apps/roles/TableHeader'
import { Box, Card, Menu, MenuItem } from '@mui/material'
import { ScannerType } from 'src/types/apps/scannerTypes'

import { formatDate, formatDateWithTime } from 'src/@core/utils/format'
import ConfirmModal from 'src/pages/modals/confirmScanner'

interface ScannerStatusType {
  [key: string]: ThemeColor
}

interface CellType {
  row: ScannerType
}

const scannerStatusObj: ScannerStatusType = {
  active: 'success',
  inactive: 'secondary'
}

const columns: GridColDef[] = [
  {
    flex: 0.2,
    maxWidth: 180,
    // field is the actual table key
    field: 'scanner',
    headerName: 'Scanner',
    renderCell: ({ row }: CellType) => {
      const { fleet_name, client } = row
      return (
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <Box sx={{ display: 'flex', alignItems: 'flex-start', flexDirection: 'column', ml: '5px' }}>
            <Typography
              noWrap
              variant='body2'
              //component={Link}
              // TODO: fixme
              //href='/apps/user/view/overview/'
              sx={{
                fontWeight: 600,
                color: 'text.primary',
                textDecoration: 'none',
                '&:hover': { color: 'primary.main' }
              }}
            >
              {fleet_name}
            </Typography>
            <Typography noWrap variant='caption'>
              {client}
            </Typography>
          </Box>
        </Box>
      )
    }
  },
  {
    flex: 0.2,
    maxWidth: 180,
    field: 'ipRange',
    headerName: 'IP Range',
    renderCell: ({ row }: CellType) => {
      const ipRange = `${row?.start_ip} - ${row?.end_ip}`

      return (
        <Typography variant='body2' noWrap>
          {ipRange}
        </Typography>
      )
    }
  },

  {
    flex: 0.15,
    maxWidth: 120,
    headerName: 'minerType',
    field: 'Miner Type',
    renderCell: ({ row }: CellType) => {
      const { miner_type } = row

      return (
        <Typography noWrap sx={{ textTransform: 'capitalize' }}>
          {miner_type}
        </Typography>
      )
    }
  },
  {
    flex: 0.1,
    maxWidth: 130,
    field: 'interval',
    headerName: 'Interval (mins)',
    renderCell: ({ row }: CellType) => {
      const { mins } = row
      return (
        <Typography noWrap sx={{ textTransform: 'capitalize' }}>
          {mins}
        </Typography>
      )
    }
  },
  {
    flex: 0.1,
    minWidth: 40,
    field: 'status',
    headerName: 'Status',
    renderCell: ({ row }: CellType) => {
      const { is_active } = row
      const status = is_active ? 'active' : 'inactive'

      return (
        <CustomChip
          skin='light'
          size='small'
          label={row.is_active ? 'Active' : 'Inactive'}
          color={scannerStatusObj[status]}
          sx={{ textTransform: 'capitalize' }}
        />
      )
    }
  },
  {
    flex: 0.35,
    minWidth: 40,
    field: 'lastUpdated',
    headerName: 'Last Updated',
    renderCell: ({ row }: CellType) => {
      const { last_run } = row

      return (
        <Typography noWrap sx={{ textTransform: 'capitalize' }}>
          {formatDateWithTime(last_run)}
        </Typography>
      )
    }
  },
  {
    flex: 0.1,
    minWidth: 100,
    sortable: false,
    field: 'actions',
    headerName: '',
    renderCell: ({ row }: CellType) => (
      <Box>
        <RowOptions row={row} />
      </Box>
    )
  }
]

const ScannerList = () => {
  //const [minerType, setMinerType] = useState<string>('')

  const [plan, setPlan] = useState<string>('')
  const [value, setValue] = useState<string>('')
  const [paginationModel, setPaginationModel] = useState({ page: 0, pageSize: 10 })

  const isDataLoading = status === 'pending'
  const isError = status === 'failed'

  const handleFilter = useCallback((val: string) => {
    setValue(val)
  }, [])

  const handlePlanChange = useCallback((e: SelectChangeEvent) => {
    setPlan(e.target.value)
  }, [])

  return (
    <Grid container spacing={6}>
      <Grid item xs={12}>
        <Card>
          {/* 
             toast notification UI goes here with store?.status conditional rendering 
          */}
          {status === 'failed' && (
            <Typography variant='body1' sx={{ mt: 3 }}>
              No results found
            </Typography>
          )}
          <TableHeader plan={plan} value={value} handleFilter={handleFilter} handlePlanChange={handlePlanChange} />
          {/* <DataGrid
            autoHeight
            rows={data}
            columns={columns}
            disableRowSelectionOnClick
            pageSizeOptions={[10, 25, 50]}
            paginationModel={paginationModel}
            onPaginationModelChange={setPaginationModel}
          /> */}
        </Card>
      </Grid>
    </Grid>
  )
}

export default ScannerList

const RowOptions = ({ row }: CellType) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)

  // Modal component related states
  const [show, setShow] = useState(false) // for modal
  const [currentAction, setCurrentAction] = useState<string | null>(null)

  const rowOptionsOpen = Boolean(anchorEl)

  const handleRowOptionsClick = (event: MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleRowOptionsClose = () => {
    setAnchorEl(null)
  }

  // TODO
  // [ ] tested
  const handleView = () => {}

  const handleEdit = () => {
    setCurrentAction('edit')
    setShow(true)
  }

  const handleRestart = () => {
    if (row?.miner_type === 'Antminer') {
      setCurrentAction('restart')
      setShow(true)
    }
  }

  const handleStop = () => {
    setCurrentAction('stop')
    setShow(true)
  }

  const handleDelete = () => {
    setCurrentAction('delete')
    setShow(true)
  }

  const onConfirmAction = () => {
    switch (currentAction) {
      case 'edit':
        break
      case 'restart':
        break
      case 'stop':
        break
      case 'delete':
        break
      default:
        break
    }
  }

  let confirmMessage
  switch (currentAction) {
    case 'restart':
      confirmMessage = `You are about to restart ${row?.fleet_name}.`
      break
    case 'stop':
      confirmMessage = `You are about to stop ${row?.fleet_name}.`
      break
    case 'delete':
      confirmMessage = `You are about to delete ${row?.fleet_name}.`
      break
    default:
      break
  }

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
        <MenuItem sx={{ '& svg': { mr: 2 } }} onClick={handleView}>
          <Icon icon='mdi:eye-outline' fontSize={20} />
          View
        </MenuItem>
        {/*
             Edit miner function should invoke 
             the modal component
             Reuse the AddUserDrawer component
        */}
        {/*
        TODO: add edit modal
        <MenuItem disabled={row?.is_active} onClick={handleRowOptionsClose} sx={{ '& svg': { mr: 2 } }}>
          <Icon icon='mdi:pencil-outline' fontSize={20} />
          Edit
        </MenuItem> */}
        <MenuItem disabled={row?.is_active} onClick={handleRestart} sx={{ '& svg': { mr: 2 } }}>
          <Icon icon='mdi:restart' fontSize={20} />
          Restart
        </MenuItem>
        <MenuItem disabled={!row?.is_active} onClick={handleStop} sx={{ '& svg': { mr: 2 } }}>
          <Icon icon='mdi:pause-circle-outline' fontSize={20} />
          Stop
        </MenuItem>
        <MenuItem disabled={row?.is_active} onClick={handleDelete} sx={{ '& svg': { mr: 2 } }}>
          <Icon icon='mdi:wifi-remove' fontSize={20} />
          Delete
        </MenuItem>
      </Menu>
      <ConfirmModal
        title={row?.fleet_name}
        message={confirmMessage}
        onConfirm={onConfirmAction}
        setShow={setShow}
        show={show}
      />
    </>
  )
}
