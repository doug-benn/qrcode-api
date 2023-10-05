import io

import qrcode
import uvicorn
from fastapi import FastAPI
from starlette.responses import StreamingResponse


def genorateqrcode():
    qr = qrcode.QRCode(
        version=1,
        error_correction=qrcode.constants.ERROR_CORRECT_L,
        box_size=10,
        border=4,
    )
    qr.add_data("Some data")
    qr.make(fit=True)

    img = qr.make_image(fill_color="black", back_color="white")
    img_bytes = io.BytesIO()
    img.save(img_bytes)
    img_bytes.seek(0)
    return img_bytes


app = FastAPI()


@app.get("/qrcode")
async def image_endpoint():
    qrcode = genorateqrcode()
    return StreamingResponse(qrcode, media_type="image/png")


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
