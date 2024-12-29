import io

import qrcode
import uvicorn
from fastapi import FastAPI
from starlette.responses import StreamingResponse


def genorateqrcode(data, version, box_size, border):
    qr = qrcode.QRCode(
        version=version,
        error_correction=qrcode.constants.ERROR_CORRECT_L,
        box_size=box_size,
        border=border,
    )
    qr.add_data(data)
    qr.make(fit=True)

    img = qr.make_image(
        fill_color="black",
        back_color="white",
    )
    img_bytes = io.BytesIO()
    img.save(img_bytes)
    img_bytes.seek(0)
    return img_bytes


app = FastAPI()


@app.get("/qrcode")
async def qrcode_endpoint(data: str = "Sample", version: int = None, box_size: int = 10, border: int = 4):
    print(data)
    qrcode = genorateqrcode(data, version, box_size, border)
    return StreamingResponse(qrcode, media_type="image/png")


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=7878)
