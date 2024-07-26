import { ChangeEvent, useState } from 'react'

import Grid from '@mui/material/Grid'

import { CustomRadioIconsData, CustomRadioIconsProps } from 'src/@core/components/custom-radio/types'

import CustomRadioIcons from 'src/@core/components/custom-radio/icons'

interface IconType {
  icon: CustomRadioIconsProps['icon']
  iconProps: CustomRadioIconsProps['iconProps']
}

const data: CustomRadioIconsData[] = [
  {
    value: 'starter',
    title: 'Starter',
    isSelected: true,
    content: 'A simple start for everyone.'
  },
  {
    value: 'standard',
    title: 'Standard',
    content: 'For small to medium businesses.'
  },
  {
    value: 'enterprise',
    title: 'Enterprise',
    content: 'Solution for big organizations.'
  }
]

const icons: IconType[] = [
  { icon: 'mdi:rocket-launch-outline', iconProps: { fontSize: '2rem', style: { marginBottom: 8 } } },
  { icon: 'mdi:account-outline', iconProps: { fontSize: '2rem', style: { marginBottom: 8 } } },
  { icon: 'mdi:crown-outline', iconProps: { fontSize: '2rem', style: { marginBottom: 8 } } }
]

const CustomRadioWithIcons = () => {
  const initialSelected: string = data.filter(item => item.isSelected)[data.filter(item => item.isSelected).length - 1]
    .value

  const [selected, setSelected] = useState<string>(initialSelected)

  const handleChange = (prop: string | ChangeEvent<HTMLInputElement>) => {
    if (typeof prop === 'string') {
      setSelected(prop)
    } else {
      setSelected((prop.target as HTMLInputElement).value)
    }
  }

  return (
    <Grid container spacing={4}>
      {data.map((item, index) => (
        <CustomRadioIcons
          key={index}
          data={data[index]}
          selected={selected}
          icon={icons[index].icon}
          name='custom-radios-icons'
          handleChange={handleChange}
          gridProps={{ sm: 4, xs: 12 }}
          iconProps={icons[index].iconProps as any}
        />
      ))}
    </Grid>
  )
}

export default CustomRadioWithIcons
