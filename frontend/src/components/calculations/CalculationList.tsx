import { useEffect, useState } from 'react'
import { api, type Calculation, type Formular } from '@/lib/api/client'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'

export function CalculationList() {
  const [calculations, setCalculations] = useState<Calculation[]>([])
  const [formulars, setFormulars] = useState<Formular[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedCalculation, setSelectedCalculation] = useState<Calculation | null>(null)
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [isFormularDialogOpen, setIsFormularDialogOpen] = useState(false)

  useEffect(() => {
    loadCalculations()
    loadFormulars()
  }, [])

  async function loadCalculations() {
    try {
      const response = await api.calculations.list()
      setCalculations(response.data)
      setError(null)
    } catch (err) {
      setError('Failed to load calculations')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  async function loadFormulars() {
    try {
      const response = await api.formulars.list()
      setFormulars(response.data)
    } catch (err) {
      console.error('Failed to load formulars:', err)
    }
  }

  async function handleCreateCalculation(name: string) {
    try {
      await api.calculations.create({ name })
      await loadCalculations()
      setIsDialogOpen(false)
    } catch (err) {
      setError('Failed to create calculation')
      console.error(err)
    }
  }

  async function handleUpdateCalculation(id: string, name: string) {
    try {
      await api.calculations.update(id, { name })
      await loadCalculations()
      setSelectedCalculation(null)
      setIsDialogOpen(false)
    } catch (err) {
      setError('Failed to update calculation')
      console.error(err)
    }
  }

  async function handleDeleteCalculation(id: string) {
    try {
      await api.calculations.delete(id)
      await loadCalculations()
    } catch (err) {
      setError('Failed to delete calculation')
      console.error(err)
    }
  }

  async function handleAddFormular(calculationId: string, formularId: string) {
    try {
      await api.calculations.addFormular(calculationId, { formularId })
      await loadCalculations()
      setIsFormularDialogOpen(false)
    } catch (err) {
      setError('Failed to add formular')
      console.error(err)
    }
  }

  async function handleRemoveFormular(calculationId: string, formularId: string) {
    try {
      await api.calculations.removeFormular(calculationId, formularId)
      await loadCalculations()
    } catch (err) {
      setError('Failed to remove formular')
      console.error(err)
    }
  }

  if (loading) return <div>Loading...</div>
  if (error) return <div className="text-red-500">{error}</div>

  return (
    <div className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Calculations</h2>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>Create Calculation</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                {selectedCalculation ? 'Edit Calculation' : 'Create Calculation'}
              </DialogTitle>
            </DialogHeader>
            <CalculationForm
              calculation={selectedCalculation}
              onSubmit={(name) =>
                selectedCalculation
                  ? handleUpdateCalculation(selectedCalculation.id, name)
                  : handleCreateCalculation(name)
              }
            />
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid gap-4">
        {calculations.map((calculation) => (
          <div
            key={calculation.id}
            className="p-4 border rounded-lg"
          >
            <div className="flex justify-between items-center mb-4">
              <h3 className="font-semibold">{calculation.name}</h3>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  onClick={() => {
                    setSelectedCalculation(calculation)
                    setIsDialogOpen(true)
                  }}
                >
                  Edit
                </Button>
                <Dialog open={isFormularDialogOpen} onOpenChange={setIsFormularDialogOpen}>
                  <DialogTrigger asChild>
                    <Button variant="outline">Add Formular</Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Add Formular to Calculation</DialogTitle>
                    </DialogHeader>
                    <FormularSelector
                      formulars={formulars}
                      onSelect={(formularId) => handleAddFormular(calculation.id, formularId)}
                    />
                  </DialogContent>
                </Dialog>
                <Button
                  variant="destructive"
                  onClick={() => handleDeleteCalculation(calculation.id)}
                >
                  Delete
                </Button>
              </div>
            </div>
            
            <div className="space-y-2">
              <h4 className="text-sm font-medium text-gray-500">Formulars:</h4>
              {calculation.formulars.map((calculationFormular) => (
                <div
                  key={calculationFormular.id}
                  className="flex justify-between items-center p-2 bg-gray-50 rounded"
                >
                  <span>{calculationFormular.formular.name}</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() =>
                      handleRemoveFormular(calculation.id, calculationFormular.formularId)
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

interface CalculationFormProps {
  calculation?: Calculation | null
  onSubmit: (name: string) => void
}

function CalculationForm({ calculation, onSubmit }: CalculationFormProps) {
  const [name, setName] = useState(calculation?.name ?? '')

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
      <Button type="submit">{calculation ? 'Update' : 'Create'}</Button>
    </form>
  )
}

interface FormularSelectorProps {
  formulars: Formular[]
  onSelect: (formularId: string) => void
}

function FormularSelector({ formulars, onSelect }: FormularSelectorProps) {
  return (
    <div className="grid gap-2">
      {formulars.map((formular) => (
        <Button
          key={formular.id}
          variant="outline"
          className="justify-start"
          onClick={() => onSelect(formular.id)}
        >
          {formular.name}
        </Button>
      ))}
    </div>
  )
}
