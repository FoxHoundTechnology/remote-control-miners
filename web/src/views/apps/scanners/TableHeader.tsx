import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import MenuItem from '@mui/material/MenuItem'
import TextField from '@mui/material/TextField'
import InputLabel from '@mui/material/InputLabel'
import FormControl from '@mui/material/FormControl'
import Select, { SelectChangeEvent } from '@mui/material/Select'

import Icon from 'src/@core/components/icon'

interface TableHeaderProps {
  client: string
  status: string
  value: string
  clientList: string[]
  handleFilter: (val: string) => void
  handleClientChange: (e: SelectChangeEvent) => void
  handleScannerStatusChange: (e: SelectChangeEvent) => void
}

const TableHeader = (props: TableHeaderProps) => {
  const { client, clientList, status, handleScannerStatusChange, handleClientChange, handleFilter, value } = props

  return (
    <Box sx={{ p: 5, pb: 3, display: 'flex', flexWrap: 'wrap', alignItems: 'center', justifyContent: 'flex-end' }}>
      <Box sx={{ display: 'flex', flexWrap: 'wrap', alignItems: 'center' }}>
        <TextField
          size='small'
          value={value}
          placeholder='Search'
          sx={{ mr: 4, mb: 2 }}
          onChange={e => handleFilter(e.target.value)}
        />
        <FormControl size='small' sx={{ mb: 2, mr: 4 }}>
          <InputLabel id='client-select'>Select Client</InputLabel>
          <Select
            size='small'
            value={client}
            id='select-client'
            label='Select Client'
            labelId='client-select'
            onChange={handleClientChange}
            inputProps={{ placeholder: 'Select Client' }}
          >
            <MenuItem value=''>Select Client</MenuItem>
            {clientList &&
              clientList.map((clientName, index) => (
                <MenuItem key={index} value={clientName}>
                  {clientName}
                </MenuItem>
              ))}
          </Select>
        </FormControl>
        <FormControl size='small' sx={{ mb: 2 }}>
          <InputLabel id='status-select'>Select Status</InputLabel>
          <Select
            size='small'
            value={status}
            id='select-status'
            label='Select Status'
            labelId='status-select'
            onChange={handleScannerStatusChange}
            inputProps={{ placeholder: 'Select Plan' }}
          >
            <MenuItem value=''>Select Status</MenuItem>
            <MenuItem value='active'>Active</MenuItem>
            <MenuItem value='inactive'>Inactive</MenuItem>
          </Select>
        </FormControl>
      </Box>
    </Box>
  )
}

export default TableHeader
