import { useEffect, useState } from 'react'
import { api, type Node } from '@/lib/api/client'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'

export function NodeList() {
  const [nodes, setNodes] = useState<Node[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedNode, setSelectedNode] = useState<Node | null>(null)
  const [isDialogOpen, setIsDialogOpen] = useState(false)

  useEffect(() => {
    loadNodes()
  }, [])

  async function loadNodes() {
    try {
      const response = await api.nodes.list()
      setNodes(response.data)
      setError(null)
    } catch (err) {
      setError('Failed to load nodes')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  async function handleCreateNode(name: string, nodeData: string) {
    try {
      await api.nodes.create({ name, nodeData })
      await loadNodes()
      setIsDialogOpen(false)
    } catch (err) {
      setError('Failed to create node')
      console.error(err)
    }
  }

  async function handleUpdateNode(id: string, name: string, nodeData: string) {
    try {
      await api.nodes.update(id, { name, nodeData })
      await loadNodes()
      setSelectedNode(null)
      setIsDialogOpen(false)
    } catch (err) {
      setError('Failed to update node')
      console.error(err)
    }
  }

  async function handleDeleteNode(id: string) {
    try {
      await api.nodes.delete(id)
      await loadNodes()
    } catch (err) {
      setError('Failed to delete node')
      console.error(err)
    }
  }

  if (loading) return <div>Loading...</div>
  if (error) return <div className="text-red-500">{error}</div>

  return (
    <div className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Nodes</h2>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>Create Node</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                {selectedNode ? 'Edit Node' : 'Create Node'}
              </DialogTitle>
            </DialogHeader>
            <NodeForm
              node={selectedNode}
              onSubmit={(name, nodeData) =>
                selectedNode
                  ? handleUpdateNode(selectedNode.id, name, nodeData)
                  : handleCreateNode(name, nodeData)
              }
            />
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid gap-4">
        {nodes.map((node) => (
          <div
            key={node.id}
            className="p-4 border rounded-lg flex justify-between items-center"
          >
            <div>
              <h3 className="font-semibold">{node.name}</h3>
              <pre className="mt-2 text-sm bg-gray-50 p-2 rounded">
                {node.nodeData}
              </pre>
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                onClick={() => {
                  setSelectedNode(node)
                  setIsDialogOpen(true)
                }}
              >
                Edit
              </Button>
              <Button
                variant="destructive"
                onClick={() => handleDeleteNode(node.id)}
              >
                Delete
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

interface NodeFormProps {
  node?: Node | null
  onSubmit: (name: string, nodeData: string) => void
}

function NodeForm({ node, onSubmit }: NodeFormProps) {
  const [name, setName] = useState(node?.name ?? '')
  const [nodeData, setNodeData] = useState(node?.nodeData ?? '')

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()
        onSubmit(name, nodeData)
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
      <div>
        <label className="block text-sm font-medium mb-1">Node Data</label>
        <textarea
          value={nodeData}
          onChange={(e) => setNodeData(e.target.value)}
          className="w-full p-2 border rounded h-32"
          required
        />
      </div>
      <Button type="submit">{node ? 'Update' : 'Create'}</Button>
    </form>
  )
}
