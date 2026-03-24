import type { MakefileDoc } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/makefile/models.js'

// 客户端序列化器，与后端 printer.go 逻辑保持一致
export function serializeDoc(doc: MakefileDoc): string {
  const lines: string[] = []

  // 1. Variables
  for (const v of doc.variables) {
    lines.push(`${v.name} ${v.operator} ${v.value}`)
  }

  // 2. .PHONY declaration
  const phonyNames = doc.targets.filter(t => t.isPhony).map(t => t.name)
  if (phonyNames.length > 0) {
    if (doc.variables.length > 0) lines.push('')
    lines.push(`.PHONY: ${phonyNames.join(' ')}`)
  }

  // 3. Targets
  if (doc.targets.length > 0) lines.push('')
  doc.targets.forEach((target, i) => {
    if (i > 0) lines.push('')
    const depPart = target.deps.length > 0 ? ` ${target.deps.join(' ')}` : ''
    lines.push(`${target.name}:${depPart}`)
    for (const cmd of target.commands) {
      lines.push(`\t${cmd}`)
    }
  })

  // 4. RawBlocks appended at the end
  for (const rb of doc.rawBlocks) {
    const content = rb.content.endsWith('\n') ? rb.content.slice(0, -1) : rb.content
    lines.push(content)
  }

  return lines.join('\n') + (lines.length > 0 ? '\n' : '')
}
