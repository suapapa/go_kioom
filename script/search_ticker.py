#!/usr/bin/env python3
# -*- coding: utf-8 -*-

# /// script
# requires-python = ">=3.9"
# dependencies = [
#     "finance-datareader",
#     "pandas",
# ]
# ///

"""
KRX Stock Ticker Search Script
Queries KRX stocks using FinanceDataReader and searches for tickers by name.
Supports exact and partial matching, quiet mode for shell scripting, and rich terminal output.

Usage:
    python3 script/search_ticker.py <stock_name>
    python3 script/search_ticker.py 삼성전자 --quiet
"""

import sys
import argparse
import unicodedata

# Gracefully handle missing dependencies and provide clear instructions
try:
    import FinanceDataReader as fdr
    import pandas as pd
except ImportError as e:
    missing_module = str(e).split("'")[-2] if "'" in str(e) else str(e)
    print("\n" + "="*80, file=sys.stderr)
    print(f"❌ Error: Missing required package '{missing_module}'", file=sys.stderr)
    print("Please install the required dependencies using the command below:", file=sys.stderr)
    print("\n    pip install finance-datareader pandas", file=sys.stderr)
    print("="*80 + "\n", file=sys.stderr)
    sys.exit(1)

def get_visual_width(s):
    """Calculate the visual width of a string in the terminal, accounting for double-width East Asian characters."""
    width = 0
    for char in s:
        if unicodedata.east_asian_width(char) in ('W', 'F', 'A'):
            width += 2
        else:
            width += 1
    return width

def pad_string(s, width, align='left'):
    """Pad a string to a visual width, taking East Asian character widths into account."""
    visual_w = get_visual_width(s)
    padding = max(0, width - visual_w)
    if align == 'left':
        return s + ' ' * padding
    elif align == 'right':
        return ' ' * padding + s
    else:
        left_pad = padding // 2
        right_pad = padding - left_pad
        return ' ' * left_pad + s + ' ' * right_pad

def format_marcap(val):
    """Format market capitalization into a human-readable Korean currency format (조/억 원)."""
    if pd.isna(val) or val == 0:
        return "-"
    # 1조 (Trillion) = 1,000,000,000,000
    # 1억 (100 Million) = 100,000,000
    trillions = val / 1_000_000_000_000
    if trillions >= 1.0:
        return f"{trillions:,.2f}조 원"
    billions = val / 100_000_000
    return f"{billions:,.1f}억 원"

def main():
    parser = argparse.ArgumentParser(description="Search KRX stock tickers by name using FinanceDataReader.")
    parser.add_argument("query", help="Name or part of the name of the stock to search.")
    parser.add_argument("-q", "--quiet", action="store_true", help="Print only the ticker code of the best match (useful for scripts).")
    args = parser.parse_args()

    query = args.query.strip()
    if not query:
        print("❌ Error: Query cannot be empty.", file=sys.stderr)
        sys.exit(1)

    # Quiet mode should suppress non-essential output to stdout, but errors should go to stderr.
    if not args.quiet:
        print("🔍 Fetching KRX stock list using FinanceDataReader...")
    
    try:
        # Fetch KRX stock list
        df_krx = fdr.StockListing('KRX')
    except Exception as e:
        print(f"❌ Error fetching KRX stock listing: {e}", file=sys.stderr)
        sys.exit(1)

    # Prepare for search (case-insensitive)
    df_krx['Name_lower'] = df_krx['Name'].astype(str).str.lower()
    query_lower = query.lower()

    # 1. Look for exact match
    exact_match = df_krx[df_krx['Name_lower'] == query_lower]
    
    if not exact_match.empty:
        results = exact_match
    else:
        # 2. Look for partial substring matches
        results = df_krx[df_krx['Name_lower'].str.contains(query_lower, na=False)]

    if results.empty:
        if not args.quiet:
            print(f"❌ No matching stocks found for '{query}'.", file=sys.stderr)
        sys.exit(1)

    # Sort results by market capitalization descending so that the most relevant/major stocks appear first
    if 'Marcap' in results.columns:
        results = results.sort_values(by='Marcap', ascending=False)

    if args.quiet:
        # Quiet mode: print ONLY the code of the top match to stdout.
        # This allows chaining, e.g., ticker=$(python3 script/search_ticker.py 삼성전자 -q)
        top_match_code = results.iloc[0]['Code']
        print(top_match_code)
        sys.exit(0)

    # Check if stdout is a TTY for color formatting
    is_tty = sys.stdout.isatty()
    COLOR_RESET = "\033[0m" if is_tty else ""
    COLOR_BOLD = "\033[1m" if is_tty else ""
    COLOR_RED = "\033[31m" if is_tty else ""      # Rise / Bullish
    COLOR_BLUE = "\033[34m" if is_tty else ""     # Fall / Bearish
    COLOR_GREEN = "\033[32m" if is_tty else ""    # Accent / Exact Match
    COLOR_GRAY = "\033[90m" if is_tty else ""

    # Detailed console formatting
    print(f"\n✨ Found {len(results)} matching stock(s) for '{query}':\n")

    # Define column layouts: name, header, and alignment
    # Code (8), Name (20), Market (8), Close (12), Change (10), Marcap (15)
    header_code = pad_string("Code", 8, 'left')
    header_name = pad_string("Name", 20, 'left')
    header_market = pad_string("Market", 8, 'left')
    header_close = pad_string("Close (KRW)", 12, 'right')
    header_change = pad_string("Change (%)", 10, 'right')
    header_marcap = pad_string("Market Cap", 15, 'right')

    header_line = f"{header_code} | {header_name} | {header_market} | {header_close} | {header_change} | {header_marcap}"
    print(COLOR_BOLD + header_line + COLOR_RESET)
    print(COLOR_GRAY + "-" * len(header_line) + COLOR_RESET)

    for _, row in results.iterrows():
        code = row.get('Code', '-')
        name = row.get('Name', '-')
        market = row.get('Market', '-')
        close = row.get('Close', 0)
        change_ratio = row.get('ChagesRatio', 0.0) # Column name spelling 'ChagesRatio' in FDR DataFrame
        marcap = row.get('Marcap', 0)

        # Format close price
        try:
            close_val = float(close)
            close_str = f"{close_val:,.0f}"
        except (ValueError, TypeError):
            close_str = str(close)

        # Format change ratio with colors
        try:
            ratio_val = float(change_ratio)
            if ratio_val > 0:
                ratio_str = f"{COLOR_RED}+{ratio_val:.2f}%{COLOR_RESET}"
                ratio_raw = f"+{ratio_val:.2f}%"
            elif ratio_val < 0:
                ratio_str = f"{COLOR_BLUE}{ratio_val:.2f}%{COLOR_RESET}"
                ratio_raw = f"{ratio_val:.2f}%"
            else:
                ratio_str = " 0.00%"
                ratio_raw = " 0.00%"
        except (ValueError, TypeError):
            ratio_str = str(change_ratio)
            ratio_raw = str(change_ratio)

        marcap_str = format_marcap(marcap)

        # Highlight exact match rows
        row_color = COLOR_GREEN if name.lower() == query_lower else ""

        # Format each column using custom Korean-safe padding
        col_code = f"{row_color}{pad_string(code, 8, 'left')}{COLOR_RESET}"
        col_name = f"{row_color}{pad_string(name, 20, 'left')}{COLOR_RESET}"
        col_market = pad_string(str(market), 8, 'left')
        col_close = pad_string(close_str, 12, 'right')
        
        # Change column length calculation needs to exclude ANSI escape sequences
        # pad_string expects clean strings to align them properly
        pad_len = 10 - len(ratio_raw)
        col_change = ' ' * max(0, pad_len) + ratio_str
        
        col_marcap = pad_string(marcap_str, 15, 'right')

        print(f"{col_code} | {col_name} | {col_market} | {col_close} | {col_change} | {col_marcap}")
    print()

if __name__ == "__main__":
    main()
