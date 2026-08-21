/**
 * ESC/POS Helper Utility for WebUSB / WebSerial Direct Thermal Printing
 */

export interface ESCPOSTicketData {
  branchName: string;
  serviceName: string;
  ticketNumber: string;
  estimatedWaitMinutes: number;
  peopleAhead: number;
  publicToken: string;
  dateStr: string;
  headerText?: string;
  footerText?: string;
  paperSize?: '58mm' | '80mm';
}

/**
 * Builds raw ESC/POS byte buffer for thermal receipt printers.
 */
export function generateESCPOSBuffer(data: ESCPOSTicketData): Uint8Array {
  const encoder = new TextEncoder();
  const bytes: number[] = [];

  // Helper push
  const pushBytes = (...items: number[]) => bytes.push(...items);
  const pushText = (text: string) => {
    const encoded = encoder.encode(text);
    encoded.forEach((b) => bytes.push(b));
  };

  // Initialize printer (ESC @)
  pushBytes(0x1b, 0x40);

  // Align Center (ESC a 1)
  pushBytes(0x1b, 0x61, 0x01);

  // Optional Custom Header
  if (data.headerText && data.headerText.trim()) {
    pushBytes(0x1b, 0x45, 0x01); // Bold ON
    pushText(data.headerText.trim() + '\n');
    pushBytes(0x1b, 0x45, 0x00); // Bold OFF
  }

  // Branch Name (Double height & width)
  pushBytes(0x1d, 0x21, 0x11); // GS ! 0x11 (2x height, 2x width)
  pushBytes(0x1b, 0x45, 0x01); // Bold ON
  pushText(data.branchName + '\n');
  pushBytes(0x1d, 0x21, 0x00); // Reset size

  // Divider
  pushText('--------------------------------\n');

  // Service Name
  pushBytes(0x1b, 0x45, 0x01); // Bold ON
  pushText(`LAYANAN: ${data.serviceName.toUpperCase()}\n`);
  pushBytes(0x1b, 0x45, 0x00);

  pushText('\nNOMOR ANTREAN:\n');

  // Big Ticket Number (4x height & width)
  pushBytes(0x1d, 0x21, 0x33); // GS ! 0x33 (4x height, 4x width)
  pushBytes(0x1b, 0x45, 0x01);
  pushText(`${data.ticketNumber}\n`);
  pushBytes(0x1d, 0x21, 0x00); // Reset size
  pushBytes(0x1b, 0x45, 0x00);

  pushText('\n');

  // Queue Stats (Left align or centered)
  pushText(`Sisa Antrean  : ${data.peopleAhead} Orang\n`);
  pushText(`Est. Waktu    : ~${data.estimatedWaitMinutes} Menit\n`);
  pushText(`Waktu Cetak   : ${data.dateStr}\n`);

  // Divider
  pushText('--------------------------------\n');

  // Custom Footer
  if (data.footerText && data.footerText.trim()) {
    pushText(data.footerText.trim() + '\n');
  } else {
    pushText('Terima kasih atas kunjungan Anda.\nHarap menunggu hingga dipanggil.\n');
  }

  // Feed 4 lines (ESC d 4)
  pushBytes(0x1b, 0x64, 0x04);

  // Cut Paper (GS V 66 0)
  pushBytes(0x1d, 0x56, 0x42, 0x00);

  return new Uint8Array(bytes);
}

/**
 * Attempts to print directly via WebUSB API if supported by the browser.
 */
export async function printDirectWebUSB(data: ESCPOSTicketData): Promise<boolean> {
  if (!('usb' in navigator)) {
    console.warn('WebUSB is not supported in this browser environment');
    return false;
  }

  try {
    const navUSB = (navigator as any).usb;
    const device = await navUSB.requestDevice({
      filters: [{ classCode: 7 }], // Printer class
    });

    await device.open();
    if (device.configuration === null) {
      await device.selectConfiguration(1);
    }
    await device.claimInterface(0);

    const buffer = generateESCPOSBuffer(data);
    // Find bulk OUT endpoint
    const endpoint = device.configuration.interfaces[0].alternate.endpoints.find(
      (e: any) => e.direction === 'out'
    );

    if (endpoint) {
      await device.transferOut(endpoint.endpointNumber, buffer);
    }
    await device.close();
    return true;
  } catch (err) {
    console.error('WebUSB Print Failed:', err);
    return false;
  }
}
