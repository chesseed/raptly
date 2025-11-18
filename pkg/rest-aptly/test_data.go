package aptly

type testPkgData struct {
	JSON string
	Pkgs []Package
}

var testPkgsDetailed = testPkgData{
	JSON: `[
    {
        "Architecture": "amd64",
        "Depends": "libc6 (>= 2.34)",
        "Description": " John's hello package\n John's package is written in C\n and prints a greeting.\n .\n It is awesome.\n",
        "Filename": "hello_3.0.0-2_amd64.deb",
        "FilesHash": "96e8a0deaf8fc95f",
        "Installed-Size": "23",
        "Key": "Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f",
        "MD5sum": "be7cbf8cf38633a26b73c4511b2d597e",
        "Maintainer": "John Doe <john@doe.com>",
        "Package": "hello",
        "Priority": "optional",
        "SHA1": "3a4c46b150d3cbe8adb27c44b5b12cca3fd63668",
        "SHA256": "52417f0e39865af616b69514bb475a2b79d3c06b02d965236e3a1e66a035cc72",
        "SHA512": "a0fc5403436286c64a8e55a885d5ca1b0ac43407550ad19a012b9cecbcae14327a8d42975672cf7f6f957e2ae812dd2159862ae143beb7982bdf698a0109bade",
        "Section": "devel",
        "ShortKey": "Pamd64 hello 3.0.0-2",
        "Size": "2648",
        "Version": "3.0.0-2"
    },
    {
        "Architecture": "amd64",
        "Auto-Built-Package": "debug-symbols",
        "Build-Ids": "7a50c209d451f1dd8c2103771fc96c2142ee059c",
        "Depends": "hello (= 3.0.0-2)",
        "Description": " debug symbols for hello\n",
        "Filename": "hello-dbgsym_3.0.0-2_amd64.deb",
        "FilesHash": "185cc47ca86a934c",
        "Installed-Size": "16",
        "Key": "Pamd64 hello-dbgsym 3.0.0-2 185cc47ca86a934c",
        "MD5sum": "1464a3c2ad70765dbc349fc4a4b6eb2a",
        "Maintainer": "John Doe <john@doe.com>",
        "Package": "hello-dbgsym",
        "Priority": "optional",
        "SHA1": "3183e2c73091e5fa992e64b8ed392a59d7442a6a",
        "SHA256": "21dc7e8f5fafcf4683c233e715860fbf38328b376f3aba8b20a70ab2843b18a8",
        "SHA512": "a01b4d7559683cf5ca752842659acd719f17fe33ece94b773ca8aa3ee9c66085899e44050b0ebf79d9a8593548d9dd6f929a55e86bc4ea6b72ec52b2a43ef9bb",
        "Section": "debug",
        "ShortKey": "Pamd64 hello-dbgsym 3.0.0-2",
        "Size": "2628",
        "Source": "hello",
        "Version": "3.0.0-2"
    },
    {
        "Architecture": "any",
        "Binary": "hello",
        "Build-Depends": "build-essential, debhelper (>= 9)",
        "Checksums-Sha1": " 3f0a502de585a30e24d7c7141559602eced32858 470 hello_3.0.0-2.dsc\n 062e2e42233c6fbe058a44e3c50ef1bf454acc96 3448 hello_3.0.0-2.tar.gz\n",
        "Checksums-Sha256": " f3767c240a5221e6122e1e561bba81ab36891218a6f5471b8705e2913df9e93c 470 hello_3.0.0-2.dsc\n b84597204d5ee78dbdc9e2fe041d93aa19c444d145e21ec16bfb4602ecb36f99 3448 hello_3.0.0-2.tar.gz\n",
        "Checksums-Sha512": " 37c9da0f380303329908d00fe0c9806b215e12721faae8e6c056a3c1f0916679800f660f51ba990ca3577303a3dd982c6900959b40052afc5c88d696ee607ab2 470 hello_3.0.0-2.dsc\n caaa02e2bc9de1d7cbfdd6c7759c974c72ec0b58650e12ad34c5b7f895e67e7d4327ce4e3256e7cfcd14ee4a306ccc3f1bd5d9bf61cedf88edbfd40e7bb59243 3448 hello_3.0.0-2.tar.gz\n",
        "Files": " 58e1956baa409b0980474b33cb5a9e99 470 hello_3.0.0-2.dsc\n 30be0886385224b34c96853cf52262fe 3448 hello_3.0.0-2.tar.gz\n",
        "FilesHash": "571d33f41765ddba",
        "Format": "1.0",
        "Key": "Psource hello 3.0.0-2 571d33f41765ddba",
        "Maintainer": "John Doe <john@doe.com>",
        "Package": "hello",
        "Package-List": " hello deb devel optional arch=any\n",
        "ShortKey": "Psource hello 3.0.0-2",
        "Version": "3.0.0-2"
    }
]`,
	Pkgs: []Package{
		{
			Architecture: "amd64",
			Key:          "Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f",
			ShortKey:     "Pamd64 hello 3.0.0-2",
			FilesHash:    "96e8a0deaf8fc95f",
			Version:      "3.0.0-2",
			Package:      "hello",
		},
		{
			Architecture: "amd64",
			Key:          "Pamd64 hello-dbgsym 3.0.0-2 185cc47ca86a934c",
			ShortKey:     "Pamd64 hello-dbgsym 3.0.0-2",
			FilesHash:    "185cc47ca86a934c",
			Version:      "3.0.0-2",
			Package:      "hello-dbgsym",
			Source:       ptr("hello"),
		},
		{
			Architecture: "any",
			Key:          "Psource hello 3.0.0-2 571d33f41765ddba",
			ShortKey:     "Psource hello 3.0.0-2",
			FilesHash:    "571d33f41765ddba",
			Version:      "3.0.0-2",
			Package:      "hello",
		},
	},
}

var testPkgsSimple1 = testPkgData{
	JSON: `["Pamd64 nano 7.2-1+deb12u1 c5d2ac1639544e75", "Psource hello 3.0.0-2 571d33f41765ddba"]`,
	Pkgs: []Package{
		{
			Key:          "Pamd64 nano 7.2-1+deb12u1 c5d2ac1639544e75",
			Architecture: "amd64",
			FilesHash:    "c5d2ac1639544e75",
			Version:      "7.2-1+deb12u1",
			Package:      "nano",
		},
		{
			Key:          "Psource hello 3.0.0-2 571d33f41765ddba",
			Architecture: "source",
			FilesHash:    "571d33f41765ddba",
			Version:      "3.0.0-2",
			Package:      "hello",
		},
	},
}

var testPkgsSimple2 = testPkgData{
	JSON: `["Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f", "Pamd64 hello-dbgsym 3.0.0-2 185cc47ca86a934c"]`,
	Pkgs: []Package{
		{
			Key:          "Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f",
			Architecture: "amd64",
			FilesHash:    "96e8a0deaf8fc95f",
			Version:      "3.0.0-2",
			Package:      "hello",
		},
		{
			Key:          "Pamd64 hello-dbgsym 3.0.0-2 185cc47ca86a934c",
			Architecture: "amd64",
			FilesHash:    "185cc47ca86a934c",
			Version:      "3.0.0-2",
			Package:      "hello-dbgsym",
		},
	},
}
