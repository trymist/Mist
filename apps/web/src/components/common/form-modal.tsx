import { useState, useEffect } from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"

interface Field {
  name: string
  label: string
  type: "text" | "email" | "password" | "textarea" | "select" | "tags"
  required?: boolean
  options?: { label: string; value: string }[]
  defaultValue?: string | string[] | number
}

interface FormModalProps<T extends Record<string, unknown>> {
  isOpen: boolean
  onClose: () => void
  title: string
  fields: Field[]
  onSubmit: (data: T) => void | Promise<void>
}

export function FormModal<T extends Record<string, unknown>>({
  isOpen,
  onClose,
  title,
  fields,
  onSubmit,
}: FormModalProps<T>) {
  const initialState = Object.fromEntries(fields.map((f) => [f.name, f.defaultValue ?? ""]))
  const [formData, setFormData] = useState<Record<string, string | string[] | number>>(initialState)
  const [tagInput, setTagInput] = useState("")
  const [tags, setTags] = useState<string[]>([])

  useEffect(() => {
    if (isOpen) {
      const state = Object.fromEntries(fields.map((f) => [f.name, f.defaultValue ?? ""]))
      setFormData(state)

      const tagsField = fields.find((f) => f.type === "tags")
      if (tagsField && Array.isArray(tagsField.defaultValue)) {
        setTags(tagsField.defaultValue)
      } else {
        setTags([])
      }
    }
  }, [isOpen, fields])

  const handleAddTag = () => {
    const newTag = tagInput.trim()
    if (!newTag) return
    if (tags.includes(newTag)) {
      toast.error("Tag already exists")
      return
    }
    setTags([...tags, newTag])
    setTagInput("")
  }

  const handleRemoveTag = (tag: string) => setTags(tags.filter((t) => t !== tag))

  const handleChange = (name: string, value: string) => {
    setFormData((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const finalData = { ...formData, tags }
    await onSubmit(finalData as unknown as T)
    onClose()
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          {fields.map((field) => {
            if (field.type === "textarea")
              return (
                <div key={field.name}>
                  <label className="text-sm text-muted-foreground">{field.label}</label>
                  <Textarea
                    value={formData[field.name]}
                    onChange={(e) => handleChange(field.name, e.target.value)}
                    required={field.required}
                    className="mt-1"
                  />
                </div>
              )

            if (field.type === "select")
              return (
                <div key={field.name}>
                  <label className="text-sm text-muted-foreground">{field.label}</label>
                  <select
                    value={formData[field.name]}
                    onChange={(e) => handleChange(field.name, e.target.value)}
                    className="w-full bg-background border rounded-md mt-1 px-3 py-2"
                  >
                    {field.options?.map((opt) => (
                      <option key={opt.value} value={opt.value}>
                        {opt.label}
                      </option>
                    ))}
                  </select>
                </div>
              )

            if (field.type === "tags")
              return (
                <div key="tags">
                  <label className="text-sm text-muted-foreground">{field.label}</label>
                  <div className="flex gap-2 mt-1">
                    <Input
                      value={tagInput}
                      onChange={(e) => setTagInput(e.target.value)}
                      onKeyDown={(e) => {
                        if (e.key === "Enter") {
                          e.preventDefault()
                          handleAddTag()
                        }
                      }}
                      placeholder="Add a tag"
                    />
                    <Button type="button" variant="secondary" onClick={handleAddTag}>
                      Add
                    </Button>
                  </div>
                  <div className="flex flex-wrap gap-2 mt-2">
                    {tags.map((tag) => (
                      <span
                        key={tag}
                        className="px-2 py-1 text-sm rounded-full bg-primary/10 text-primary flex items-center gap-2"
                      >
                        {tag}
                        <button
                          type="button"
                          onClick={() => handleRemoveTag(tag)}
                          className="text-destructive hover:underline"
                        >
                          ×
                        </button>
                      </span>
                    ))}
                  </div>
                </div>
              )

            return (
              <div key={field.name}>
                <label className="text-sm text-muted-foreground">{field.label}</label>
                <Input
                  type={field.type}
                  value={formData[field.name]}
                  onChange={(e) => handleChange(field.name, e.target.value)}
                  required={field.required}
                  className="mt-1"
                />
              </div>
            )
          })}

          <DialogFooter className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit">Submit</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
