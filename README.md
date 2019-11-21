# mock

New()
↓
Mocker.Mock(tags string, data interface{})
↓isPtr
mock(tag Tag, v reflect.Value)
↓Ptr        ↓Struct            ↓Slice      ↓Array       ↓Map            ↓Field
mock        mockStruct         mockSlice   mockArray    mockMap         mockField
            ↓canSet&&notIgnore ↓setLen     ↓            ↓isMapKeyString ↓
            ↓mock              mock        mock         ↓setMapSize     mockString
                                                        mockString      mockInt
                                                        mock            mockUint
                                                                        mockFloat
                                                        