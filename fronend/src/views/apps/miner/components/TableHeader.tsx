import { useState } from 'react'

import { FormControl, InputAdornment, MenuItem } from '@mui/material'
import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import Select, { SelectChangeEvent } from '@mui/material/Select'
import TextField from '@mui/material/TextField'

import SearchIcon from '@mui/icons-material/Search'
import SortIcon from '@mui/icons-material/Sort'

import Icon from 'src/@core/components/icon'
import { AppDispatch } from 'src/store'
import ConfirmModal from 'src/pages/modals/confirm'
import { Command } from 'src/store/apps/minerControl'

import { FilterValue } from '../../../../pages/apps/miner/list'

interface TableHeaderProps {
  filters: FilterValue
  initialFilters: FilterValue
  setFilters: React.Dispatch<React.SetStateAction<FilterValue>>
  remoteControlCallback: (val: number[], command: Command) => void
  toggleFilter: () => void
  openFilter: boolean
  rowSelectionModel: any
  setRowSelectionModel: any
}
/*
  TODO:
  [ ] Add the toggle button for remote control 
  [ ] miner mode selector as well 
*/

const TableHeader = (props: TableHeaderProps) => {
  const { filters, setFilters, initialFilters, remoteControlCallback, openFilter, toggleFilter, rowSelectionModel } =
    props
  const [remoteControlSelected, setRemoteControlSelected] = useState(false)
  const [show, setShow] = useState(false)
  const [command, setCommand] = useState<Command>()

  const remoteControlButtonHandler = () => {
    setRemoteControlSelected(!remoteControlSelected)
  }

  const handleFilterChange =
    (filterKey: keyof FilterValue) => (event: SelectChangeEvent | React.ChangeEvent<HTMLInputElement>) => {
      const value = event.target.value
      setFilters(prev => ({
        ...prev,
        [filterKey]: value === 'reset' ? initialFilters[filterKey] : value
      }))
    }

  const messageGenerator = () => {
    if (command === Command.LowPower) {
      const num = rowSelectionModel.length
      return `You are about to set ${num} miners to Low Power Mode.`
    } else if (command === Command.Sleep) {
      const num = rowSelectionModel.length
      return `You are about to set ${num} miners to Sleep Mode.`
    } else if (command === Command.Reboot) {
      const num = rowSelectionModel.length
      return `You are about to Reboot ${num} miners.`
    } else {
      // normal mode
      const num = rowSelectionModel.length
      return `You are about to set ${num} miners to Normal Mode.`
    }
  }

  return (
    <Box
      sx={{
        p: { xs: 2, sm: 4 },
        pb: 3,
        display: 'flex',
        flexWrap: 'wrap',
        alignItems: 'center',
        justifyContent: 'space-between'
      }}
    >
      <Box sx={{ mb: { xs: 2, sm: 0 } }}>
        <Button
          size='large'
          sx={{ mr: 2, mb: 1, pr: 3, pl: 3 }}
          color='secondary'
          variant='outlined'
          startIcon={<SortIcon />}
          onClick={toggleFilter}
        >
          {openFilter ? 'Reset Filter' : 'Filter'}
        </Button>
        <Button
          sx={{ mr: 2, mb: 1, pr: 3, pl: 3 }}
          size='large'
          color='primary'
          variant='outlined'
          startIcon={<Icon icon='mdi:shape-circle-plus' fontSize={20} />}
          onClick={remoteControlButtonHandler}
        >
          Remote Control
        </Button>
        {remoteControlSelected && (
          <>
            <Select size='small' sx={{ mr: 4, mb: 1, pr: 4, pl: 4, pt: 1 }}>
              <MenuItem value='normal' onClick={() => setCommand(Command.Normal)}>
                Normal
              </MenuItem>
              <MenuItem value='lowpower' onClick={() => setCommand(Command.LowPower)}>
                Low Power
              </MenuItem>
              <MenuItem value='sleep' onClick={() => setCommand(Command.Sleep)}>
                Sleep
              </MenuItem>
              <MenuItem value='reboot' onClick={() => setCommand(Command.Reboot)}>
                Reboot
              </MenuItem>
              {/* TODO: Bulk Config is Under Development */}
              <MenuItem disabled value='config'>
                Configuration
              </MenuItem>
            </Select>
            <Button
              sx={{ mr: 4, mb: 1, pr: 4, pl: 4 }}
              onClick={() => {
                setShow(true)
              }}
              variant='contained'
              // TODO : button active or not
              disabled={rowSelectionModel.length === 0}
            >
              Submit
            </Button>
          </>
        )}
      </Box>
      <Box sx={{ mb: { xs: 2, sm: 0 } }}>
        <Box sx={{ display: 'flex', flexWrap: 'wrap', alignItems: 'center' }}>
          <Button
            size='medium'
            sx={{ mr: 2, pr: 3, pl: 3 }}
            color='secondary'
            variant='outlined'
            startIcon={<Icon icon='mdi:export-variant' fontSize={25} />}
          >
            Export
          </Button>
          <TextField
            size='small'
            value={filters.search}
            InputProps={{
              startAdornment: (
                <InputAdornment position='start'>
                  <SearchIcon />
                </InputAdornment>
              )
            }}
            placeholder='Search'
            onChange={event => handleFilterChange('search')(event as React.ChangeEvent<HTMLInputElement>)}
          />
        </Box>
      </Box>
      <ConfirmModal
        show={show}
        setShow={setShow}
        message={messageGenerator()}
        onConfirm={() => remoteControlCallback(rowSelectionModel, command)}
        setRowSelectionModel={props?.setRowSelectionModel}
      />
    </Box>
  )
}

export default TableHeader
