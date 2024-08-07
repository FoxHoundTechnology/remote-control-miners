import React, { useState, useCallback } from 'react'
import { Box, Button, TextField, Typography, Paper } from '@mui/material'
import { useDrag, useDrop, XYCoord } from 'react-dnd'

const ContainerType = 'container'

// Define the grid size (in pixels)
const GRID_SIZE = 20

// Helper function to snap a value to the grid
const snapToGrid = (value: number): number => {
  return Math.round(value / GRID_SIZE) * GRID_SIZE
}

interface ContainerData {
  id: string
  name: string
  description: string
  left: number
  top: number
}

interface DraggableContainerProps {
  container: ContainerData
  moveContainer: (id: string, left: number, top: number) => void
  deleteContainer: (id: string) => void
}

interface DragItem {
  id: string
  left: number
  top: number
}

const DraggableContainer: React.FC<DraggableContainerProps> = ({ container, moveContainer, deleteContainer }) => {
  const [{ isDragging }, drag] = useDrag(
    () => ({
      type: ContainerType,
      item: { id: container.id, left: container.left, top: container.top },
      collect: monitor => ({
        isDragging: monitor.isDragging()
      })
    }),
    [container.id, container.left, container.top]
  )

  return (
    <Paper
      ref={drag}
      style={{
        position: 'absolute',
        left: container.left,
        top: container.top,
        opacity: isDragging ? 0.5 : 1,
        cursor: 'move'
      }}
      elevation={3}
      sx={{
        padding: 2,
        backgroundColor: '#f0f0f0',
        width: 200
      }}
    >
      <Typography variant='h6'>{container.name}</Typography>
      <Typography variant='body2'>{container.description}</Typography>
      <Button onClick={() => deleteContainer(container.id)} size='small' color='error'>
        Delete
      </Button>
    </Paper>
  )
}

const DataCenterSitemapContent: React.FC = () => {
  const [containers, setContainers] = useState<ContainerData[]>([])
  const [newContainer, setNewContainer] = useState<{ name: string; description: string }>({ name: '', description: '' })

  const moveContainer = useCallback((id: string, left: number, top: number) => {
    setContainers(prevContainers =>
      prevContainers.map(container =>
        container.id === id ? { ...container, left: snapToGrid(left), top: snapToGrid(top) } : container
      )
    )
  }, [])

  const [, drop] = useDrop(
    () => ({
      accept: ContainerType,
      drop(item: DragItem, monitor) {
        const delta = monitor.getDifferenceFromInitialOffset() as XYCoord
        const left = Math.round(item.left + delta.x)
        const top = Math.round(item.top + delta.y)
        moveContainer(item.id, left, top)
        return undefined
      }
    }),
    [moveContainer]
  )

  const handleAddContainer = () => {
    if (newContainer.name && newContainer.description) {
      const newContainerData: ContainerData = {
        ...newContainer,
        id: Date.now().toString(),
        left: snapToGrid(Math.random() * (window.innerWidth - 200)),
        top: snapToGrid(Math.random() * (window.innerHeight - 200))
      }
      setContainers(prevContainers => [...prevContainers, newContainerData])
      setNewContainer({ name: '', description: '' })
    }
  }

  const handleDeleteContainer = (id: string) => {
    setContainers(prevContainers => prevContainers.filter(container => container.id !== id))
  }

  return (
    <Box sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ padding: 2, borderBottom: '1px solid #ccc' }}>
        <Typography variant='h4' gutterBottom>
          Data Center Sitemap
        </Typography>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <TextField
            label='Container Name'
            value={newContainer.name}
            onChange={e => setNewContainer(prev => ({ ...prev, name: e.target.value }))}
            size='small'
          />
          <TextField
            label='Container Description'
            value={newContainer.description}
            onChange={e => setNewContainer(prev => ({ ...prev, description: e.target.value }))}
            size='small'
          />
          <Button variant='contained' onClick={handleAddContainer}>
            Add Container
          </Button>
        </Box>
      </Box>

      <Box
        ref={drop}
        sx={{
          flexGrow: 1,
          position: 'relative',
          overflow: 'hidden'
        }}
      >
        {containers.map(container => (
          <DraggableContainer
            key={container.id}
            container={container}
            moveContainer={moveContainer}
            deleteContainer={handleDeleteContainer}
          />
        ))}
      </Box>
    </Box>
  )
}

export default DataCenterSitemapContent
