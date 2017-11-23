#!/usr/bin/env python3

import gzip
import decimal
import random
import sys
import math

ops = {
    "*": "multiplication",
    "+": "addition",
    "-": "subtraction",
    "/": "division",
    "qC": "comparison",
    "quant": "quantize",
    "A": "abs",
    "cfd": "convert-to-string",
    "~": "neg",
    "*-": "fused-multiply-add",

    # Custom
    "rat": "convert-to-rat",
    "sign": "sign",
    "signbit": "signbit",
}

modes = {
    "=0": decimal.ROUND_HALF_EVEN,
    "=^": decimal.ROUND_HALF_UP,
    "0": decimal.ROUND_DOWN,
    "<": decimal.ROUND_FLOOR,
    ">": decimal.ROUND_CEILING,
    # decimal.ROUND_HALF_DOWN,
    "^": decimal.ROUND_UP,
    # decimal.ROUND_05UP,
}


def rand_bool():
    return random.randint(0, 1) % 2 == 0


def make_dec():
    sign = "+" if rand_bool() else "-"
    if rand_bool():
        f = random.uniform(1, sys.float_info.max)
    else:
        f = random.getrandbits(5000)

    r = random.randint(0, 25)
    if r == 0:
        f = math.nan
    elif r == 1:
        f = math.inf
    return decimal.Decimal("{}{}".format(sign, f))


DEC_TEN = decimal.Decimal(10)


def rand_dec(quant=False):
    with decimal.localcontext() as ctx:
        ctx.clear_traps()

        if quant:
            x = random.randint(0, 500)
            if rand_bool():
                d = DEC_TEN ** x
            else:
                d = DEC_TEN ** -x
        else:
            d = make_dec()
            for _ in range(random.randint(0, 3)):
                q = random.randint(0, 4)
                if q == 1:
                    d *= make_dec()
                elif q == 2:
                    d /= make_dec()
                elif q == 3:
                    d -= make_dec()
                elif q == 4:
                    d += make_dec()
                # else: == 0
    return d


def conv(x):
    if x is None:
        return None
    if isinstance(x, decimal.Decimal):
        if x.is_infinite():
            return '-Inf' if x.is_signed() else 'Inf'
        else:
            return x
    return x


def write_line(out, prec, op, mode, r, x, y=None, z0=None, flags=None):
    if x is None:
        raise ValueError("bad args")

    x = conv(x)
    y = conv(y)
    z0 = conv(z0)
    r = conv(r)
    if y is not None:
        if z0 is not None:
            str = "d{}{} {} {} {} {} -> {} {}\n".format(
                prec, op, mode, x, y, z0, r, flags)
        else:
            str = "d{}{} {} {} {} -> {} {}\n".format(
                prec, op, mode, x, y, r, flags)
    else:
        str = "d{}{} {} {} -> {} {}\n".format(
            prec, op, mode, x, r, flags)
    out.write(str)


def perform_op(op):
    r = None
    x = rand_dec()
    y = None  # possibly unused
    z0 = None  # possibly unused

    try:
        # Binary
        if op == "*":
            y = rand_dec()
            r = x * y
        elif op == "+":
            y = rand_dec()
            r = x + y
        elif op == "-":
            y = rand_dec()
            r = x - y
        elif op == "/":
            y = rand_dec()
            r = x / y
        elif op == "qC":
            y = rand_dec()
            r = x.compare(y)
        elif op == "quant":
            y = rand_dec(True)
            r = x.quantize(y)
            y = -y.as_tuple().exponent

        # Unary
        elif op == "A":
            r = decimal.getcontext().abs(x)
        elif op == "cfd":
            r = str(x)
        elif op == "rat":
            while True:
                try:
                    n, d = x.as_integer_ratio()
                    break
                except Exception:  # ValueError if nan, etc.
                    x = rand_dec()
            x, y = n, d
            # 'suite' can't parse the following:
            # r = "{}/{}".format(n, d)
            r = "#"
        elif op == "sign":
            if x < 0:
                r = -1
            elif x > 0:
                r = +1
            else:
                r = 0
        elif op == "signbit":
            r = x.is_signed()
        elif op == "~":
            r = -x

        # Ternary
        elif op == "*-":
            y = rand_dec()
            z0 = rand_dec()
            r = x.fma(y, z0)
        else:
            raise ValueError("bad op {}".format(op))
    except Exception as e:
        if op == "quant":
            raise e
        pass

    return (r, x, y, z0)


traps = {
    decimal.Clamped: "c",
    decimal.DivisionByZero: "z",
    decimal.Inexact: "x",
    decimal.InvalidOperation: "i",
    decimal.Overflow: "o",
    decimal.Rounded: "r",
    decimal.Subnormal: "s",
    decimal.Underflow: "u",
    decimal.FloatOperation: "***",
}


def rand_traps():
    t = {}
    s = ""
    for key, val in traps.items():
        b = key != decimal.FloatOperation and rand_bool()
        if b:
            s += val
        t[key] = int(b)
    return (t, s)

    # set N higher for local testing.
N = 1000


def make_tables():
    for op, name in ops.items():
        with gzip.open("{}-tables.gz".format(name), "wt") as f:
            for i in range(1, N):
                mode = random.choice(list(modes.keys()))
                # t, ts = rand_traps()

                ctx = decimal.getcontext()
                ctx.Emax = decimal.MAX_EMAX
                ctx.Emin = decimal.MIN_EMIN
                ctx.rounding = modes[mode]
                ctx.prec = random.randint(1, 5000)
                ctx.clear_traps()
                ctx.clear_flags()

                r, x, y, z0 = perform_op(op)

                conds = ""
                for key, value in ctx.flags.items():
                    if value == 1 and key != decimal.FloatOperation:
                        conds += traps[key]
                write_line(f, ctx.prec, op, mode, r, x, y, z0, conds)


make_tables()
