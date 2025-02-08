import { useEffect, useState } from 'react'
import { api, type Formular, type Node } from '@/lib/api/client'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'

export function FormularList() {
  const [formulars, setFormulars] = useState<Formular[]>([])
  const [nodes, setNodes] = useState<Node[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedFormular, setSelectedFormular] = useState<Formular | null>(null)
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [isNodeDialogOpen, setIsNodeDialogOpen] = useState(false)

  useEffect(() => {
    loadFormulars()
    loadNodes()
  }, [])

  async function loadFormulars() {
    try {
      const response = await api.formulars.list()
      setFormulars(response.data)
      setError(null)
    } catch (err) {
      setError('Failed to load formulars')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  async function loadNodes() {
    try {
      const response = await api.nodes.list()
      setNodes(response.data)
    } catch (err) {
      console.error('Failed to load nodes:', err)
    }
  }

  async function handleCreateFormular(name: string) {
    try {
      await api.formulars.create({ name })
      await loadFormulars()
      setIsDialogOpen(false)
    } catch (err) {
      setError('Failed to create formular')
      console.error(err)
    }
  }

  async function handleUpdateFormular(id: string, name: string) {
    try {
      await api.formulars.update(id, { name })
      await loadFormulars()
      setSelectedFormular(null)
      setIsDialogOpen(false)
    } catch (err) {
      setError('Failed to update formular')
      console.error(err)
    }
  }

  async function handleDeleteFormular(id: string) {
    try {
      await api.formulars.delete(id)
      await loadFormulars()
    } catch (err) {
      setError('Failed to delete formular')
      console.error(err)
    }
  }

  async function handleAddNode(formularId: string, nodeId: string) {
    try {
      await api.formulars.addNode(formularId, { nodeId })
      await loadFormulars()
      setIsNodeDialogOpen(false)
    } catch (err) {
      setError('Failed to add node')
      console.error(err)
    }
  }

  async function handleRemoveNode(formularId: string, nodeId: string) {
    try {
      await api.formulars.removeNode(formularId, nodeId)
      await loadFormulars()
    } catch (err) {
      setError('Failed to remove node')
      console.error(err)
    }
  }

  if (loading) return <div>Loading...</div>
  if (error) return <div className="text-red-500">{error}</div>

  return (
    <div className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Formulars</h2>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>Create Formular</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                {selectedFormular ? 'Edit Formular' : 'Create Formular'}
              </DialogTitle>
            </DialogHeader>
            <FormularForm
              formular={selectedFormular}
              onSubmit={(name) =>
                selectedFormular
                  ? handleUpdateFormular(selectedFormular.id, name)
                  : handleCreateFormular(name)
              }
            />
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid gap-4">
        {formulars.map((formular) => (
          <div
            key={formular.id}
            className="p-4 border rounded-lg"
          >
            <div className="flex justify-between items-center mb-4">
              <h3 className="font-semibold">{formular.name}</h3>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  onClick={() => {
                    setSelectedFormular(formular)
                    setIsDialogOpen(true)
                  }}
                >
                  Edit
                </Button>
                <Dialog open={isNodeDialogOpen} onOpenChange={setIsNodeDialogOpen}>
                  <DialogTrigger asChild>
                    <Button variant="outline">Add Node</Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Add Node to Formular</DialogTitle>
                    </DialogHeader>
                    <NodeSelector
                      nodes={nodes}
                      onSelect={(nodeId) => handleAddNode(formular.id, nodeId)}
                    />
                  </DialogContent>
                </Dialog>
                <Button
                  variant="destructive"
                  onClick={() => handleDeleteFormular(formular.id)}
                >
                  Delete
                </Button>
              </div>
            </div>
            
            <div className="space-y-2">
              <h4 className="text-sm font-medium text-gray-500">Nodes:</h4>
              {formular.nodes.map((formularNode) => (
                <div
                  key={formularNode.id}
                  className="flex justify-between items-center p-2 bg-gray-50 rounded"
                >
                  <span>{formularNode.node.name}</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() =>
                      handleRemoveNode(formular.id, formularNode.nodeId)
                    }
                  >
                    Remove
                  </Button>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

interface FormularFormProps {
  formular?: Formular | null
  onSubmit: (name: string) => void
}

function FormularForm({ formular, onSubmit }: FormularFormProps) {
  const [name, setName] = useState(formular?.name ?? '')

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()
        onSubmit(name)
      }}
      className="space-y-4"
    >
      <div>
        <label className="block text-sm font-medium mb-1">Name</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="w-full p-2 border rounded"
          required
        />
      </div>
      <Button type="submit">{formular ? 'Update' : 'Create'}</Button>
    </form>
  )
}

interface NodeSelectorProps {
  nodes: Node[]
  onSelect: (nodeId: string) => void
}

function NodeSelector({ nodes, onSelect }: NodeSelectorProps) {
  return (
    <div className="grid gap-2">
      {nodes.map((node) => (
        <Button
          key={node.id}
          variant="outline"
          className="justify-start"
          onClick={() => onSelect(node.id)}
        >
          {node.name}
        </Button>
      ))}
    </div>
  )
}
