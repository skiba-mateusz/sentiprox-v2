export const FONT_MONO = '"IBM Plex Mono", monospace';
export const ACCENT = "#F0B429";

export const CATEGORY_COLORS = {
  benign: "#3FB68B",
  sqli: "#E8634A",
  xss: "#B57BE8",
  cmdi: "#E85AA0",
  traversal: "#4A9FE0",
  "brute-force": "#EB5E6A",
  dos: "#F0885E",
  recon: "#6E7BE0",
  enumeration: "#45B8B0",
  blocked: "#97A2B8",
};

export function catColor(c: string) {
  return (CATEGORY_COLORS as Record<string, string>)[c] || "#97A2B8";
}

export const tooltipStyle = {
  background: "#141A28",
  border: "1px solid #2A3346",
  borderRadius: 8,
  fontFamily: FONT_MONO,
  fontSize: 12,
  color: "#E6EAF2",
};
